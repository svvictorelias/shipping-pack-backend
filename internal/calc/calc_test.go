package calc

import "testing"

func TestCalculatePacksEdgeCase(t *testing.T) {
	packs := []int{23, 31, 53}
	target := 500000
	counts, total, packCount, err := CalculatePacks(target, packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 500000 {
		t.Fatalf("expected total 500000 got %d", total)
	}
	expected := 2 + 7 + 9429
	if packCount != expected {
		t.Fatalf("expected packCount %d got %d", expected, packCount)
	}
	if counts[23] != 2 || counts[31] != 7 || counts[53] != 9429 {
		t.Fatalf("unexpected counts: %#v", counts)
	}
}

// Target zero deve retornar erro
func TestCalculatePacks_TargetZero(t *testing.T) {
	_, _, _, err := CalculatePacks(0, []int{10, 20})
	if err == nil {
		t.Fatalf("expected error for target=0, got nil")
	}
}

// Lista vazia de packs deve retornar erro
func TestCalculatePacks_EmptyPacks(t *testing.T) {
	_, _, _, err := CalculatePacks(100, []int{})
	if err == nil {
		t.Fatalf("expected error for empty packs")
	}
}

// Target pequeno deve usar o menor pack
func TestCalculatePacks_SmallTarget(t *testing.T) {
	packs := []int{10, 20, 50}
	target := 5
	counts, total, packCount, err := CalculatePacks(target, packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 10 {
		t.Fatalf("expected total 10 got %d", total)
	}
	if packCount != 1 {
		t.Fatalf("expected 1 pack got %d", packCount)
	}
	if counts[10] != 1 {
		t.Fatalf("expected 1 pack of 10, got %#v", counts)
	}
}

// Target impossível de atingir deve retornar erro
func TestCalculatePacks_NoSolution(t *testing.T) {
	packs := []int{4, 6}
	target := 7
	counts, total, packCount, err := CalculatePacks(target, packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total < target {
		t.Fatalf("total %d is less than target %d", total, target)
	}
	if total != 8 {
		t.Fatalf("expected total 8 got %d", total)
	}
	if packCount != 2 || counts[4] != 2 {
		t.Fatalf("expected two packs of 4, got %#v", counts)
	}
}

// Packs únicos (deve ser simples de calcular)
func TestCalculatePacks_SinglePack(t *testing.T) {
	packs := []int{100}
	target := 250
	counts, total, packCount, err := CalculatePacks(target, packs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 300 { // 3 * 100 >= 250
		t.Fatalf("expected total 300 got %d", total)
	}
	if counts[100] != 3 {
		t.Fatalf("expected 3 packs of 100, got %#v", counts)
	}
	if packCount != 3 {
		t.Fatalf("expected packCount=3 got %d", packCount)
	}
}
