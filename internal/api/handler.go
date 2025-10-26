package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
	"github.com/svvictorelias/shipping-pack-backend/internal/service"

	_ "github.com/lib/pq"
)

// Server holds dependencies for HTTP handlers.
type Server struct {
	svc *service.Service
	db  *sql.DB
}

// NewServer builds server given a store implementation.
func NewServer(svc *service.Service, db *sql.DB) *Server {
	return &Server{svc: svc, db: db}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/packs", s.packsHandler)
	mux.HandleFunc("/calculate", s.calculateHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	return c.Handler(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "ts": time.Now().Format(time.RFC3339)})
}

func (s *Server) packsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		packs, err := s.svc.GetPacks()
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"packs": packs})
		return
	case http.MethodPost:
		var body struct {
			Packs []int `json:"packs"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json")
			return
		}
		if len(body.Packs) == 0 {
			writeErr(w, http.StatusBadRequest, "packs required")
			return
		}
		if err := s.svc.SetPacks(body.Packs); err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	default:
		writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

func (s *Server) calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, http.StatusMethodNotAllowed, "method")
		return
	}
	var body struct {
		Items int `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	if body.Items <= 0 {
		writeErr(w, http.StatusBadRequest, "items must be > 0")
		return
	}
	packs, err := s.svc.GetPacks()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	counts, total, packCount, err := s.svc.Calculate(body.Items, packs)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := map[string]interface{}{
		"counts":      counts,
		"total_items": total,
		"pack_count":  packCount,
		"waste":       total - body.Items,
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// SetupDB helper to open DB based on env DATABASE_URL
func SetupDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, ErrNoDBURL
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// quick ping
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

var ErrNoDBURL = &customErr{"DATABASE_URL not set"}

type customErr struct{ s string }

func (e *customErr) Error() string { return e.s }
