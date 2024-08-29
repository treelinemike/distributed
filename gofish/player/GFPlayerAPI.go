package main

import (
	"engg415/playingcards"
	"errors"
	"fmt"
)

type GFPlayerAPI int

var hand = new(playingcards.Deck)

// wrapper for AddCard() that statisfies the RPC interface
func (gfapi *GFPlayerAPI) AddCard(card playingcards.Card, resp *int) error {
	err := hand.AddCard(card)
	if err != nil {
		return errors.New("could not add card")
	}
	fmt.Println("We got a card! Current hand:")
	hand.Show()
	return nil
}

func (gfapi *GFPlayerAPI) ResetHand(args int, resp *int) error {
	fmt.Println("Resetting deck")
	hand.Reset()
	return nil
}

func (gfapi *GFPlayerAPI) TakeTopCard(_ int, c *playingcards.Card) error {
	*c = hand.TakeTopCard()
	return nil
}
