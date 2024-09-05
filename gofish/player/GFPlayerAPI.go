package main

import (
	"engg415/gofish/gfcommon"
	"engg415/playingcards"
	"errors"
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
			log.Printf("Removed a book of value %d\n", k)
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

// RPC to give cards to another player
func (gfapi *GFPlayerAPI) GiveCards(requestedVal int, transferredCards *[]playingcards.Card) error {
	newHand := make([]playingcards.Card, 0)
	for _, c := range hand.Cards {
		if c.Val == requestedVal {
			*transferredCards = append(*transferredCards, c)
		} else {
			newHand = append(newHand, c)
		}
		hand.Cards = newHand
	}
	log.Printf("Gave up %d cards to another player!\n", len(*transferredCards))
	return nil
}

// RPC to allow a player to take a turn
// this is the main logic of the game
func (gfapi *GFPlayerAPI) TakeTurn(_ int, resp *gfcommon.GFPlayerReturn) error {

	var j int
	c := new(playingcards.Card)
	tryAgain := true // set up to continue unless we (later) set otherwise

	// handle blink(1) indicator and logging for turn start/end
	log.Printf("Starting turn with %d books and hand: %s\n", numBooks, hand.String())
	_, err := exec.Command("blink1-on.sh").Output()
	if err != nil {
		log.Println("Could not turn on blink(1) indicator")
	}
	defer func() {
		_, err = exec.Command("blink1-off.sh").Output()
		if err != nil {
			log.Println("Could not turn off blink(1) indicator")
		}
		log.Printf("Ending turn with %d books and hand: %s\n", numBooks, hand.String())
	}()

	// query players for cards until our luck runs out
	for tryAgain {

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

		// remove books
		removeBooksFromHand()

		// select a card at random from our hand
		// TODO: we could be smarter about the choice!
		// but, values repeated in hand are more likely to be chosen which is good...
		valToRequest := hand.Cards[rand.Intn(hand.NumCards())].Val

		// select a random player to request cards from
		// TODO: we could be smarter about this choice as well!
		// TODO: decide whether to fail if only one player (zero others), not failing can be helpful for debugging - play with deck only
		if len(config.OtherPlayers) > 0 {
			playerToRequestFrom := rand.Intn(len(config.OtherPlayers))
			log.Printf("Requesting card value %d from player index %d\n", valToRequest, playerToRequestFrom)

			// connect to specified player
			opponent, err := rpc.DialHTTP("tcp", config.OtherPlayers[playerToRequestFrom].Address+":"+config.OtherPlayers[playerToRequestFrom].Port)
			if err != nil {
				log.Println("Error connecting to client: " + config.OtherPlayers[playerToRequestFrom].Address)
			}
			log.Println("Connected to client")

			// try to get cards from opponent
			newCards := make([]playingcards.Card, 0)
			err = opponent.Call("GFPlayerAPI.GiveCards", valToRequest, &newCards)
			if err != nil {
				log.Fatalf("Couldn't request cards from player: %s\n", config.OtherPlayers[playerToRequestFrom].Address)
			}
			log.Printf("Received %d cards from %s\n", len(newCards), config.OtherPlayers[playerToRequestFrom].Address)

			// add cards to deck if we received any
			if len(newCards) > 0 {
				hand.Cards = append(hand.Cards, newCards...)
				removeBooksFromHand()
			} else {
				// try to pull a card from the deck
				log.Println("Attempting to pull card from deck")
				err = host.Call("GFHostAPI.TakeTopCard", j, c)
				if err != nil {
					log.Fatal("Error retrieving top card from deck")
				}
				if c.Val == 0 {
					log.Println("Deck is empty, cannot continue this turn!")
					resp.NumBooks = numBooks
					resp.NumCardsInHand = hand.NumCards()
					return nil
				} else {
					log.Printf("Pulled a card from the deck: %s\n", c.String())
					if c.Val == valToRequest {
						log.Println("You fished your wish!")
					} else {
						tryAgain = false
					}
					hand.AddCard((*c))
					log.Printf("Added card, hand is now: %s\n", hand.String())
					removeBooksFromHand()
				}
				tryAgain = false
			}
		}
	} // end for tryAgain

	// return
	resp.NumBooks = numBooks
	resp.NumCardsInHand = hand.NumCards()
	return nil
}
