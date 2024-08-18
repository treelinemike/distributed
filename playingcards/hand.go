package playingcards

// the "Hand" structre holding info about the current hand
// must have a capital "H" indicating that it is exported
type Hand struct {
	cards []Card
}

func (hand *Hand) AddCard(c Card) error { // TODO: can we leverage non-class methods in Go for this (use same method as in Deck?
    hand.cards = append(hand.cards, c)
    return nil
}


