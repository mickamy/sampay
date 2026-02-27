package model

// CalcAmounts computes the split amount for each participant based on tier
// weights. The tier value itself is the weight (e.g. tier=3 means weight 3).
// Returns a map from participant ID to the amount in yen.
// When there are no participants or total weight is zero, returns an empty map.
func CalcAmounts(totalAmount int, participants []EventParticipant) map[string]int {
	totalWeight := 0
	for _, p := range participants {
		totalWeight += p.Tier
	}
	if totalWeight == 0 {
		return map[string]int{}
	}

	amounts := make(map[string]int, len(participants))
	for _, p := range participants {
		amounts[p.ID] = totalAmount * p.Tier / totalWeight
	}
	return amounts
}
