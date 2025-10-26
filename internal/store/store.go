package store

// Store defines persistence operations used by the service.
// This allows easy mocking for tests.
type Store interface {
	// GetPacks returns configured pack sizes (unsorted possible).
	GetPacks() ([]int, error)

	// SetPacks atomically replaces pack sizes in DB.
	SetPacks([]int) error

	// SaveCalculation persists a run of CalculatePacks for auditing.
	// counts is map[packSize]quantity
	SaveCalculation(items int, totalItems int, packCount int, counts map[int]int) error
}
