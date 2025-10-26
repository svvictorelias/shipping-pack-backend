package store

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// PostgresStore implements Store using database/sql and Postgres.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore opens a DB connection. dbURL is a Postgres connection string.
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// GetPacks returns pack sizes from DB.
func (s *PostgresStore) GetPacks() ([]int, error) {
	rows, err := s.db.Query("SELECT size FROM packs ORDER BY size ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packs := []int{}
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, err
		}
		packs = append(packs, size)
	}
	return packs, rows.Err()
}

// SetPacks replaces pack sizes atomically.
func (s *PostgresStore) SetPacks(packs []int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// clear table then insert
	if _, err := tx.Exec("DELETE FROM packs"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO packs(size, created_at) VALUES($1,$2)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range packs {
		if _, err := stmt.Exec(p, time.Now().UTC()); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// SaveCalculation saves calculation summary and items.
func (s *PostgresStore) SaveCalculation(items int, totalItems int, packCount int, counts map[int]int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var calcID int
	err = tx.QueryRow(
		"INSERT INTO calculations(items,total_items,pack_count,created_at) VALUES($1,$2,$3,$4) RETURNING id",
		items, totalItems, packCount, time.Now().UTC(),
	).Scan(&calcID)
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO calculation_items(calculation_id, pack_size, quantity) VALUES($1,$2,$3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for size, qty := range counts {
		if _, err := stmt.Exec(calcID, size, qty); err != nil {
			return err
		}
	}
	return tx.Commit()
}
