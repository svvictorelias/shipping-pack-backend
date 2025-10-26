package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/svvictorelias/shipping-pack-backend/internal/api"
	"github.com/svvictorelias/shipping-pack-backend/internal/service"
	"github.com/svvictorelias/shipping-pack-backend/internal/store"

	_ "github.com/lib/pq"
)

func main() {
	// Setup DB
	db, err := api.SetupDB()
	if err != nil {
		log.Printf("DB not available: %v. Falling back to mock store (development).", err)
		// fallback to mock store to allow local dev without DB
		mock := store.NewMockStore([]int{250, 500, 1000, 2000, 5000})
		svc := service.NewService(mock)
		srv := api.NewServer(svc, nil)
		startHTTP(srv)
		return
	}

	// create Postgres store
	pstore := store.NewPostgresStore(db)
	svc := service.NewService(pstore)

	srv := api.NewServer(svc, db)
	startHTTP(srv)
}

func startHTTP(srv *api.Server) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	h := srv.Routes()
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        h,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("listening on :%s", port)
	log.Fatal(s.ListenAndServe())
}
