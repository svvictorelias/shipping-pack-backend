package store

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostgresStore_GetPacks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"size"}).AddRow(100).AddRow(200)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT size FROM packs ORDER BY size ASC")).WillReturnRows(rows)

	store := NewPostgresStore(db)
	packs, err := store.GetPacks()
	if err != nil {
		t.Fatalf("GetPacks error: %v", err)
	}
	if len(packs) != 2 || packs[0] != 100 || packs[1] != 200 {
		t.Fatalf("unexpected packs: %v", packs)
	}
}

func TestPostgresStore_SetPacks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM packs").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare("INSERT INTO packs").ExpectExec().
		WithArgs(100, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	store := NewPostgresStore(db)
	err = store.SetPacks([]int{100})
	if err != nil {
		t.Fatalf("SetPacks error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestPostgresStore_SaveCalculation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO calculations(items,total_items,pack_count,created_at) VALUES($1,$2,$3,$4) RETURNING id")).
		WithArgs(450, 500, 5, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectPrepare("INSERT INTO calculation_items").
		ExpectExec().
		WithArgs(1, 100, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	store := NewPostgresStore(db)
	err = store.SaveCalculation(450, 500, 5, map[int]int{100: 2})
	if err != nil {
		t.Fatalf("SaveCalculation error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
