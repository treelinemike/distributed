package main

import (
	"engg415/gofish/gfcommon"
	"engg415/playingcards"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os/exec"
	"time"
)

type GFPlayerAPI int

var hand = new(playingcards.Deck) // it may seem strange that 1) we don't use the int value of the GFPlayerAPI type when instantiated in gofish_player.go, and 2) that our methods on the GFPlayerAPI type are all set up to act on the unexported hand instance of playingcards.Deck
var config = new(gfcommon.GFPlayerConfig)
var host *rpc.Client

// RPC allowing host to set player configuration
func (gfapi *GFPlayerAPI) SetConfig(c gfcommon.GFPlayerConfig, resp *int) error {
	*config = c
	log.Println("Config set via RPC")

	// connect to host
	var err error
	host, err = rpc.DialHTTP("tcp", config.Host.Address+":"+config.Host.Port)
	if err != nil {
		log.Println("Error connecting to host: " + config.Host.Address)
	}
	log.Println("Connected to host")
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

	// set up to continue unless we (later) set otherwise
	tryAgain := true

	// handle blink(1) indicator and logging for turn start/end
	log.Println("Taking turn")
	_, err := exec.Command("blink1-on.sh").Output()
	if err != nil {
		log.Println("Could not turn on blink(1) indicator")
	}
	defer func() {
		_, err = exec.Command("blink1-off.sh").Output()
		if err != nil {
			log.Println("Could not turn off blink(1) indicator")
		}
		log.Printf("Turn complete, hand contains %d cards\n", hand.NumCards())
	}()

	// try to take top card
	var j int
	var cards []playingcards.Card
	err = host.Call("GFHostAPI.TakeTopCard", j, &cards)

	if err != nil {
		log.Fatal("Error retrieving top card from deck")
	}
	switch len(cards) {
	case 0:
		log.Printf("Deck is empty, no card added")
	case 1:
		// TODO: add card to deck
		log.Printf("Took top card from deck: %s\n", cards[0].String())
	default:
		log.Fatal("Received more than one card from deck")
	}

	// if hand is empty try to get a card from the deck
	if hand.NumCards() == 0 {

	}

	if tryAgain {
		fmt.Println("hello")
	}

	time.Sleep(2 * time.Second)
	*resp = hand.NumCards()
	return nil
}
