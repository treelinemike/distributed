package main

import (
	"engg415/playingcards"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type GFAPI int

var hand = new(playingcards.Deck)

// wrapper for AddCard() that statisfies the RPC interface
func (gfapi *GFAPI) AddCardRPC(card playingcards.Card, resp *int) error {
	err := hand.AddCard(card)
	if err != nil {
		return errors.New("could not add card")
	}
	fmt.Println("We got a card! Current hand:")
	hand.Show()
	return nil
}

func (gfapi *GFAPI) ResetHandRPC(args int, resp *int) error {
	fmt.Println("Resetting deck")
	hand.Reset()
	return nil
}

func (gfapi *GFAPI) TakeTopCardRPC(_ int, c *playingcards.Card) error {
	*c = hand.TakeTopCard()
	return nil
}

func main() {

	fmt.Println("Starting player")

	// create a hand as just an empty deck
	//hand := new(playingcards.Deck)
	//tc := new(playingcards.TestCard)
	//thisCard := new(playingcards.Card)
	//hand.AddCard(*thisCard)
	//hand.Show()

	thisapi := new(GFAPI)

	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)

	fmt.Println("Waiting for hand to be dealt")
	for {
	}
}
