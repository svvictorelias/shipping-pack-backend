package service

import (
	"github.com/svvictorelias/shipping-pack-backend/internal/calc"
	"github.com/svvictorelias/shipping-pack-backend/internal/store"
)

// Service holds business logic and interacts with the store.
type Service struct {
	store store.Store
}

// NewService constructs service with given store.
func NewService(s store.Store) *Service {
	return &Service{store: s}
}

// GetPacks returns pack sizes from persistence.
func (s *Service) GetPacks() ([]int, error) {
	return s.store.GetPacks()
}

// SetPacks stores new pack sizes.
func (s *Service) SetPacks(packs []int) error {
	return s.store.SetPacks(packs)
}

// Calculate performs algorithm and persists the calculation result.
func (s *Service) Calculate(items int, packs []int) (map[int]int, int, int, error) {
	counts, total, packCount, err := calc.CalculatePacks(items, packs)
	if err != nil {
		return nil, 0, 0, err
	}
	// persist result (best-effort; propagate error)
	if perr := s.store.SaveCalculation(items, total, packCount, counts); perr != nil {
		// return both results and error so caller can decide; here we return error
		return counts, total, packCount, perr
	}
	return counts, total, packCount, nil
}
