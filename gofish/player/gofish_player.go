package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

func main() {

	fmt.Println("Starting player")

	// create a hand as just an empty deck
	//hand := new(playingcards.Deck)
	//tc := new(playingcards.TestCard)
	//thisCard := new(playingcards.Card)
	//hand.AddCard(*thisCard)
	//hand.Show()

	thisapi := new(GFPlayerAPI)

	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)

	fmt.Println("Waiting for hand to be dealt")

	// TODO: Find a better way to do this!
	for {
	}
}
