package calc

import (
	"errors"
	"sort"
)

// CalculatePacks finds a combination of pack sizes that achieves total >= target
// minimizing two objectives in order:
// 1) minimal total items shipped (S >= target)
// 2) among totals with same S, minimal number of packs
//
// Returns counts map[packSize]quantity, totalItems, packCount, error.
func CalculatePacks(target int, packs []int) (map[int]int, int, int, error) {
	if target <= 0 {
		return nil, 0, 0, errors.New("target must be positive")
	}
	if len(packs) == 0 {
		return nil, 0, 0, errors.New("packs empty")
	}

	// copy and sort ascending for DP optimization
	p := make([]int, len(packs))
	copy(p, packs)
	sort.Ints(p)
	maxP := p[len(p)-1]

	// DP only needs to consider totals up to target + maxP - 1
	limit := target + maxP - 1
	const INF = int(1e9)

	dp := make([]int, limit+1)   // dp[s] = minimal number of packs to make exactly s (INF if unreachable)
	prev := make([]int, limit+1) // prev[s] = last pack size used to reach s

	for i := 1; i <= limit; i++ {
		dp[i] = INF
		prev[i] = -1
	}
	dp[0] = 0
	prev[0] = -1

	// fill DP
	for s := 1; s <= limit; s++ {
		for _, pack := range p {
			if pack > s {
				break
			}
			if dp[s-pack] != INF {
				if dp[s] > dp[s-pack]+1 {
					dp[s] = dp[s-pack] + 1
					prev[s] = pack
				}
			}
		}
	}

	// find minimal reachable total S >= target
	bestS := -1
	for s := target; s <= limit; s++ {
		if dp[s] != INF {
			bestS = s
			break
		}
	}
	if bestS == -1 {
		return nil, 0, 0, errors.New("no solution")
	}

	// reconstruct counts
	counts := make(map[int]int)
	s := bestS
	packCount := 0
	for s > 0 {
		pk := prev[s]
		if pk <= 0 {
			return nil, 0, 0, errors.New("reconstruction failed")
		}
		counts[pk]++
		packCount++
		s -= pk
		// safety
		if packCount > limit+10 {
			return nil, 0, 0, errors.New("reconstruction loop")
		}
	}

	return counts, bestS, packCount, nil
}
