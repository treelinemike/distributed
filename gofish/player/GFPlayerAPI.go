package main

import (
	"engg415/gofish/gfcommon"
	"engg415/playingcards"
	"errors"
	"log"
	"os/exec"
	"time"
)

type GFPlayerAPI int

var hand = new(playingcards.Deck) // it may seem strange that 1) we don't use the int value of the GFPlayerAPI type when instantiated in gofish_player.go, and 2) that our methods on the GFPlayerAPI type are all set up to act on the unexported hand instance of playingcards.Deck
var config = new(gfcommon.GFPlayerConfig)

func (gfapi *GFPlayerAPI) SetConfig(c gfcommon.GFPlayerConfig, resp *int) error {
	*config = c
	log.Println("Config set via RPC")
	return nil
}

// RPC wrapper for playingcards.Deck.AddCard()
func (gfapi *GFPlayerAPI) AddCardToHand(card playingcards.Card, resp *int) error {
	err := hand.AddCard(card)
	if err != nil {
		return errors.New("could not add card")
	}
	log.Println("Received a card via RPC")
	return nil
}

// RPC wrapper for playingcards.Deck.Reset()
func (gfapi *GFPlayerAPI) ResetHand(args int, resp *int) error {
	hand.Reset()
	log.Println("Deck reset via RPC")
	return nil
}

// RPC wrapper for playingcards.Deck.TakeTopCard()
func (gfapi *GFPlayerAPI) TakeTopCard(_ int, c *playingcards.Card) error {
	*c = hand.TakeTopCard()
	return nil
}

func (gfapi *GFPlayerAPI) TakeTurn(_ int, resp *int) error {
	log.Println("Taking my turn...")
	_, err := exec.Command("blink1-on.sh").Output()
	if err != nil {
		log.Println("Could not turn on blink(1) indicator")
	}
	time.Sleep(2 * time.Second)
	_, err = exec.Command("blink1-off.sh").Output()
	if err != nil {
		log.Println("Could not turn off blink(1) indicator")
	}
	return nil
}
