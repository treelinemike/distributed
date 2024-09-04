package main

import (
	"engg415/gofish/gfcommon"
	"engg415/playingcards"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"os/exec"
)

type GFPlayerAPI int

var hand = new(playingcards.Deck) // it may seem strange that 1) we don't use the int value of the GFPlayerAPI type when instantiated in gofish_player.go, and 2) that our methods on the GFPlayerAPI type are all set up to act on the unexported hand instance of playingcards.Deck
var config = new(gfcommon.GFPlayerConfig)
var host *rpc.Client
var numBooks int

func removeBooksFromHand() error {

	// create a map of (card rank) int -> (count) int
	counts := make(map[int]int)
	for _, c := range hand.Cards {
		counts[c.Val] += 1
	}

	// remove any books
	for k, v := range counts {
		if v == 4 {
			newHand := make([]playingcards.Card, 0)
			for _, thisCard := range hand.Cards {
				if thisCard.Val != k {
					newHand = append(newHand, thisCard)
				}
			}
			hand.Cards = newHand
			numBooks += 1
		}
	}

	// return
	return nil
}

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

func (gfapi *GFPlayerAPI) TakeTurn(_ int, resp *gfcommon.GFPlayerReturn) error {

	var j int
	c := new(playingcards.Card)
	tryAgain := true // set up to continue unless we (later) set otherwise

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

	// if hand is empty, try to take a card from the deck
	if hand.NumCards() == 0 {
		c = new(playingcards.Card)
		err = host.Call("GFHostAPI.TakeTopCard", j, c)
		if err != nil {
			log.Fatal("Error retrieving top card from deck")
		}
		if c.Val == 0 {
			log.Println("Deck is empty, cannot do anything this turn!")
			resp.NumBooks = numBooks
			resp.NumCardsInHand = hand.NumCards()
			return nil
		} else {
			hand.AddCard((*c))
			log.Printf("Added card, hand is now: %s\n", hand.String())
		}
	}

	// remove books (really only applicable on first turn in rare case that we're dealt a book)
	removeBooksFromHand()

	// query players for cards until our luck runs out
	for tryAgain {

		// select a card at random from our hand
		// TODO: we could be smarter about the choice!
		// but, if a value is repeated in the hand it will be more likely to be chosen, which is good...
		valToRequest := hand.Cards[rand.Intn(hand.NumCards())].Val

		// request card from a random player
		// TODO: we could be smarter about this choice as well!
		// TODO: decide whether to fail if only one player (helpful for debugging to play with deck)
		if len(config.OtherPlayers) > 0 {
			log.Printf("Requesting card value %d from player index %d\n", valToRequest, playerToRequestFrom)

		}

		fmt.Printf("Hand before books removed: %s\n", hand.String())
		removeBooksFromHand()
		fmt.Printf("Hand after books removed: %s\n", hand.String())
		//time.Sleep(5 * time.Second)
		tryAgain = false
	}

	// return
	resp.NumBooks = numBooks
	resp.NumCardsInHand = hand.NumCards()
	return nil
}
