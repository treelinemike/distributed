package main

import (
	"engg415/playingcards"
)

type GFHostAPI int

var deck = new(playingcards.Deck)

// RPC wrapper for TakeTopCard()
// TODO: put these in a queue to restore if needed?
// returns a slice of cards b/c that is easy to check for null, although we could return a card with zero value
func (gfapi *GFHostAPI) TakeTopCard(_ int, c *playingcards.Card) error {
	if deck.NumCards() > 0 {
		*c = deck.TakeTopCard()
	} else {
		c = new(playingcards.Card)
	}
	return nil
}
