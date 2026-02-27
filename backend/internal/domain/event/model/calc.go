package model

// CalcTierAmounts computes the per-tier amount and sets it on each EventTier
// in-place. Returns the remainder (totalAmount - sum of all tier amounts × counts).
// The tier value itself is the weight (e.g. tier=3 means weight 3).
// When there are no tiers or total weight is zero, returns 0.
func CalcTierAmounts(totalAmount int, tiers []EventTier) int {
	var totalWeight int
	for _, t := range tiers {
		totalWeight += t.Tier * t.Count
	}
	if totalWeight == 0 {
		return 0
	}

	sum := 0
	for i := range tiers {
		tiers[i].Amount = totalAmount * tiers[i].Tier / totalWeight
		sum += tiers[i].Amount * tiers[i].Count
	}
	return totalAmount - sum
}
