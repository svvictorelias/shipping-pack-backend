package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/svvictorelias/shipping-pack-backend/internal/service"
	"github.com/svvictorelias/shipping-pack-backend/internal/store"
)

// setupServer helper que cria o servidor HTTP com MockStore
func setupServer() *Server {
	mock := store.NewMockStore([]int{23, 31, 53})
	svc := service.NewService(mock)
	return NewServer(svc, nil)
}

func TestHealthHandler(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("expected status ok got %v", resp)
	}
}

func TestGetPacksHandler(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest(http.MethodGet, "/packs", nil)
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
	var body map[string][]int
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(body["packs"]) == 0 {
		t.Fatalf("expected packs not empty")
	}
}

func TestPostPacksHandler(t *testing.T) {
	srv := setupServer()
	payload := []byte(`{"packs":[10,20,30]}`)
	req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp map[string]bool
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if !resp["ok"] {
		t.Fatalf("expected ok=true, got %#v", resp)
	}
}

func TestCalculateHandler_Success(t *testing.T) {
	srv := setupServer()
	payload := []byte(`{"items":500000}`)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d, body=%s", rec.Code, rec.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["total_items"].(float64) < 500000 {
		t.Fatalf("expected total >= 500000 got %v", resp["total_items"])
	}
}

func TestCalculateHandler_InvalidJSON(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}

func TestCalculateHandler_MissingItems(t *testing.T) {
	srv := setupServer()
	payload := []byte(`{"items":0}`)
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
}
func TestPostPacksHandler_InvalidJSON(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewReader([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("invalid json")) {
		t.Fatalf("expected 'invalid json' message, got %s", rec.Body.String())
	}
}

func TestPostPacksHandler_MissingField(t *testing.T) {
	srv := setupServer()
	req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewReader([]byte(`{"wrong_field":[1,2,3]}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rec.Code)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("packs required")) {
		t.Fatalf("expected 'packs required' message, got %s", rec.Body.String())
	}
}

func TestPostPacksHandler_InternalError(t *testing.T) {
	// cria um mock store que falha no SetPacks
	mock := &failingStore{}
	svc := service.NewService(mock)
	srv := NewServer(svc, nil)

	payload := []byte(`{"packs":[1,2,3]}`)
	req := httptest.NewRequest(http.MethodPost, "/packs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rec.Code)
	}
}

func TestHandlers_MethodNotAllowed(t *testing.T) {
	srv := setupServer()

	// tentar PUT em /packs
	req := httptest.NewRequest(http.MethodPut, "/packs", nil)
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rec.Code)
	}

	// tentar GET em /calculate
	req2 := httptest.NewRequest(http.MethodGet, "/calculate", nil)
	rec2 := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 got %d", rec2.Code)
	}
}

// mock para simular erro interno
type failingStore struct{}

func (f *failingStore) GetPacks() ([]int, error)                         { return []int{10}, nil }
func (f *failingStore) SetPacks(packs []int) error                       { return errors.New("db fail") }
func (f *failingStore) SaveCalculation(int, int, int, map[int]int) error { return nil }
