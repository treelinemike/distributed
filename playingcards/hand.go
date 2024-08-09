package playingcards

// the "Hand" structre holding info about the current hand
// must have a capital "H" indicating that it is exported
type Hand struct {
	numCards uint // initialized to zero!
}

// method to update number of cards in hand
func (h *Hand) AddCards(numCardsToAdd uint, reply *int) error {
	h.numCards += numCardsToAdd
	return nil
}

