package main

import (
	"engg415/gofish/gfcommon"
	"engg415/playingcards"
	"errors"
	"log"
	"math/rand"
	"net/rpc"
	"os/exec"
	"time"
)

type GFPlayerAPI int

var hand = new(playingcards.Deck) // it may seem strange that 1) we don't use the int value of the GFPlayerAPI type when instantiated in gofish_player.go, and 2) that our methods on the GFPlayerAPI type are all set up to act on the unexported hand instance of playingcards.Deck
var config = new(gfcommon.GFPlayerConfig)
var host *rpc.Client
var numBooks int
var gameover bool

func removeBooksFromHand() error {

	// create a map of (card rank) int -> (count) int
	counts := make(map[int]int)
	for _, c := range hand.Cards {
		counts[c.Val] += 1
	}

	// remove any books from hand
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
			log.Printf("Removed a book of rank %s\n", playingcards.NumToCardChar(k))
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
		log.Printf("Error connecting to host %s:%s\n", config.Host.Address, config.Host.Port)
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
	numBooks = 0
	log.Println("Deck reset via RPC")
	return nil
}

// RPC wrapper for playingcards.Deck.TakeTopCard()
func (gfapi *GFPlayerAPI) TakeTopCard(_ int, c *playingcards.Card) error {
	*c = hand.TakeTopCard()
	return nil
}

// RPC to print current deck to log
// used, e.g. after dealing is complete
func (gfapi *GFPlayerAPI) PrintHand(_ int, resp *int) error {
	log.Printf("Current hand (%s)\n", hand.String())
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
	if len(*transferredCards) > 0 {
		log.Printf("Relinquished %d cards of rank %s\n", len(*transferredCards), playingcards.NumToCardChar(requestedVal))
	} else {
		log.Printf("Rejected request for cards of rank %s\n", playingcards.NumToCardChar(requestedVal))
	}
	return nil
}

// RPC to end game
// if this player lost the game winStatus will be zero
// otherwise, winStatus will reflect the total number of winners
func (gfapi *GFPlayerAPI) EndGame(winStatus int, resp *int) error {

	switch winStatus {
	case 0:
		log.Printf("Lost game with %d book(s) collected", numBooks)
	case 1:
		log.Printf("Won game with %d book(s) collected", numBooks)
		exec.Command("blink1-glimmer.sh").Output() // don't handle an error on this, ok if it fails (i.e. no blink1 configured)
	default:
		log.Printf("Tied with %d other player(s) for the win with %d books collected", winStatus-1, numBooks)
	}

	// close connection to host
	host.Close()
	gameover = true
	return nil
}

// RPC to allow a player to take a turn
// this is the main logic of the game
func (gfapi *GFPlayerAPI) TakeTurn(_ int, resp *gfcommon.GFPlayerReturn) error {

	var j int
	c := new(playingcards.Card)
	tryAgain := true // set up to continue unless we (later) set otherwise

	// handle blink(1) indicator and logging for turn start/end
	log.Printf("Starting turn with %d books and hand (%s)\n", numBooks, hand.String())
	exec.Command("blink1-on.sh").Output() // don't handle an error on this, ok if it fails (i.e. no blink1 configured)
	defer func() {
		time.Sleep(500 * time.Millisecond)     // delay so we can watch gameplay
		exec.Command("blink1-off.sh").Output() // don't handle an error on this, ok if it fails (i.e. no blink1 configured)
		log.Printf("Ending turn with %d books and hand (%s)\n", numBooks, hand.String())
	}()

	// query players for cards until our luck runs out
	for tryAgain {

		// if hand is empty, try to take a card from the deck
		if hand.NumCards() == 0 {
			c = new(playingcards.Card)
			err := host.Call("GFHostAPI.TakeTopCard", j, c)
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

		newCards := make([]playingcards.Card, 0)
		if len(config.OtherPlayers) > 0 {
			playerToRequestFrom := rand.Intn(len(config.OtherPlayers))
			log.Printf("Fishing for a card of rank %s\n", playingcards.NumToCardChar(valToRequest))
			log.Printf("Asking player %s:%s\n", config.OtherPlayers[playerToRequestFrom].Address, config.OtherPlayers[playerToRequestFrom].Port)

			// connect to specified player
			opponent, err := rpc.DialHTTP("tcp", config.OtherPlayers[playerToRequestFrom].Address+":"+config.OtherPlayers[playerToRequestFrom].Port)
			if err != nil {
				log.Printf("Error connecting to client %s:%s\n", config.OtherPlayers[playerToRequestFrom].Address, config.OtherPlayers[playerToRequestFrom].Port)
			}

			// try to get cards from opponent
			err = opponent.Call("GFPlayerAPI.GiveCards", valToRequest, &newCards)
			if err != nil {
				log.Fatalf("Couldn't request cards from player: %s:%s\n", config.OtherPlayers[playerToRequestFrom].Address, config.OtherPlayers[playerToRequestFrom].Port)
			}
			opponent.Close()
			log.Printf("Received %d cards from %s:%s\n", len(newCards), config.OtherPlayers[playerToRequestFrom].Address, config.OtherPlayers[playerToRequestFrom].Port)
		}

		// add cards to deck if we received any
		// otherwise try to draw from deck
		if len(newCards) > 0 {
			hand.Cards = append(hand.Cards, newCards...)
			removeBooksFromHand()
			log.Printf("Current hand (%s)\n", hand.String())
		} else {
			// try to pull a card from the deck
			c = new(playingcards.Card) // need to reset card b/c a zero value in struct from RPC won't get gobbed, so old value will persist!
			err := host.Call("GFHostAPI.TakeTopCard", j, c)
			if err != nil {
				log.Fatal("Error retrieving top card from deck")
			}
			if c.Val == 0 {
				log.Println("Could not pull from deck (deck is empty), cannot continue this turn!")
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
				removeBooksFromHand()
			}
			log.Printf("Current hand (%s)\n", hand.String())
		}

	} // end for tryAgain

	// return
	resp.NumBooks = numBooks
	resp.NumCardsInHand = hand.NumCards()
	return nil
}
