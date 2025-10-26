package store

import "sync"

// MockStore is a simple in-memory implementation of Store for unit tests.
type MockStore struct {
	mu           sync.RWMutex
	packs        []int
	calculations []mockCalc
}

type mockCalc struct {
	items     int
	total     int
	packCount int
	counts    map[int]int
}

// NewMockStore constructs a mock store pre-seeded with packs.
func NewMockStore(packs []int) *MockStore {
	return &MockStore{packs: packs}
}

func (m *MockStore) GetPacks() ([]int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]int, len(m.packs))
	copy(out, m.packs)
	return out, nil
}

func (m *MockStore) SetPacks(packs []int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.packs = make([]int, len(packs))
	copy(m.packs, packs)
	return nil
}

func (m *MockStore) SaveCalculation(items int, totalItems int, packCount int, counts map[int]int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cpy := make(map[int]int)
	for k, v := range counts {
		cpy[k] = v
	}
	m.calculations = append(m.calculations, mockCalc{
		items:     items,
		total:     totalItems,
		packCount: packCount,
		counts:    cpy,
	})
	return nil
}

// CountCalculations returns how many calculations have been saved.
// Exported so tests in other packages can assert persistence behavior.
func (m *MockStore) CountCalculations() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.calculations)
}

// LastCalculation returns the last saved calculation (items, total, packCount, counts, ok).
// ok == false when there is no saved calculation.
func (m *MockStore) LastCalculation() (items int, total int, packCount int, counts map[int]int, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.calculations) == 0 {
		return 0, 0, 0, nil, false
	}
	last := m.calculations[len(m.calculations)-1]
	// return a defensive copy of counts
	cpy := make(map[int]int, len(last.counts))
	for k, v := range last.counts {
		cpy[k] = v
	}
	return last.items, last.total, last.packCount, cpy, true
}
