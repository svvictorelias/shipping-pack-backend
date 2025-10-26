package service

import (
	"errors"
	"testing"

	"github.com/svvictorelias/shipping-pack-backend/internal/store"
)

func TestServiceCalculatePersists(t *testing.T) {
	mock := store.NewMockStore([]int{23, 31, 53})
	svc := NewService(mock)

	counts, total, _, err := svc.Calculate(500000, []int{23, 31, 53})
	if err != nil {
		t.Fatalf("calculate err: %v", err)
	}
	if total != 500000 {
		t.Fatalf("expected total 500000 got %d", total)
	}
	if counts[53] == 0 {
		t.Fatalf("expected some 53 packs")
	}
	if mock.CountCalculations() == 0 {
		t.Fatalf("expected calculations persisted")
	}
	itemsPersisted, totalPersisted, _, countsPersisted, ok := mock.LastCalculation()
	if !ok {
		t.Fatalf("expected last calculation to be available")
	}
	if totalPersisted != total {
		t.Fatalf("persisted total mismatch")
	}
	if itemsPersisted != 500000 {
		t.Fatalf("persisted items mismatch")
	}
	if countsPersisted[53] == 0 {
		t.Fatalf("persisted counts missing expected 53 packs")
	}
}

func TestServiceGetAndSetPacks(t *testing.T) {
	mock := store.NewMockStore([]int{100, 200})
	svc := NewService(mock)

	// GetPacks
	got, err := svc.GetPacks()
	if err != nil {
		t.Fatalf("GetPacks err: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 packs, got %d", len(got))
	}

	// SetPacks alters the state
	newPacks := []int{10, 20, 30}
	if err := svc.SetPacks(newPacks); err != nil {
		t.Fatalf("SetPacks err: %v", err)
	}
	got2, _ := svc.GetPacks()
	if len(got2) != 3 {
		t.Fatalf("expected 3 packs, got %d", len(got2))
	}
}

// mock that always returns errors
type errStore struct{}

func (e *errStore) GetPacks() ([]int, error) { return nil, errors.New("fail GetPacks") }
func (e *errStore) SetPacks([]int) error     { return errors.New("fail SetPacks") }
func (e *errStore) SaveCalculation(int, int, int, map[int]int) error {
	return errors.New("fail SaveCalculation")
}

func TestServiceCalculate_SaveFails(t *testing.T) {
	svc := NewService(&errStore{})
	// with errStore, SaveCalculation needs to fails
	_, _, _, err := svc.Calculate(100, []int{10, 20})
	if err == nil {
		t.Fatalf("expected error from SaveCalculation")
	}
}

func TestServiceCalculate_InvalidInput(t *testing.T) {
	mock := store.NewMockStore([]int{10})
	svc := NewService(mock)
	// call with target=0
	_, _, _, err := svc.Calculate(0, []int{10})
	if err == nil {
		t.Fatalf("expected error for target=0")
	}
}

func TestServiceCalculate_EmptyPacks(t *testing.T) {
	mock := store.NewMockStore([]int{})
	svc := NewService(mock)
	_, _, _, err := svc.Calculate(100, []int{})
	if err == nil {
		t.Fatalf("expected error for empty packs")
	}
}
