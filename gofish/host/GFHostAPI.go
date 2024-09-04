package main

import (
	"engg415/playingcards"
)

type GFHostAPI int

var deck = new(playingcards.Deck)

// RPC wrapper for TakeTopCard()
// TODO: put these in a queue to restore if needed?
func (gfapi *GFHostAPI) TakeTopCard(_ int, c *[]playingcards.Card) error {
	if deck.NumCards() > 0 {
		*c = append(*c, deck.TakeTopCard())
	}
	return nil
}
