package model

// CalcAmounts computes the split amount for each participant based on tier
// weights and sets the Amount field on each participant in-place.
// The tier value itself is the weight (e.g. tier=3 means weight 3).
// When there are no participants or total weight is zero, no amounts are set.
func CalcAmounts(totalAmount int, participants []EventParticipant) {
	var totalWeight int
	for _, p := range participants {
		totalWeight += p.Tier
	}
	if totalWeight == 0 {
		return
	}

	for i := range participants {
		participants[i].Amount = totalAmount * participants[i].Tier / totalWeight
	}
}
