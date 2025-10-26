package main

import (
	"log"

	"github.com/svvictorelias/shipping-pack-backend/internal/api"

	"github.com/svvictorelias/go-migrate/pkg/migrate"
)

func main() {
	// Setup DB
	db, err := api.SetupDB()
	if err != nil {
		log.Fatalf("DB not available: %v", err)
	}
	if err := migrate.Run(db, "postgres", "internal/database/migrations", false); err != nil {
		log.Fatalf("Erro ao aplicar migrations: %v", err)
	}
}
