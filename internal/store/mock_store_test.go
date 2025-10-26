package store

import "testing"

func TestMockStore_GetAndSetPacks(t *testing.T) {
	ms := NewMockStore([]int{100, 200})
	packs, err := ms.GetPacks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(packs) != 2 {
		t.Fatalf("expected 2 packs got %d", len(packs))
	}

	newPacks := []int{10, 20, 30}
	if err := ms.SetPacks(newPacks); err != nil {
		t.Fatalf("SetPacks error: %v", err)
	}

	got, _ := ms.GetPacks()
	if len(got) != 3 {
		t.Fatalf("expected 3 packs got %d", len(got))
	}
}

func TestMockStore_SaveCalculationAndCount(t *testing.T) {
	ms := NewMockStore([]int{50, 100})
	counts := map[int]int{50: 2, 100: 3}

	err := ms.SaveCalculation(450, 500, 5, counts)
	if err != nil {
		t.Fatalf("SaveCalculation error: %v", err)
	}

	if ms.CountCalculations() != 1 {
		t.Fatalf("expected 1 calculation saved, got %d", ms.CountCalculations())
	}

	items, total, packs, savedCounts, ok := ms.LastCalculation()
	if !ok {
		t.Fatal("expected last calculation to exist")
	}
	if items != 450 || total != 500 || packs != 5 {
		t.Fatalf("unexpected saved data: %d %d %d", items, total, packs)
	}
	if savedCounts[100] != 3 {
		t.Fatalf("expected 3 of 100, got %v", savedCounts)
	}
}
