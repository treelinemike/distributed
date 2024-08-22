package main

import (
	"engg415/playingcards"
	"fmt"
    "net"
    "net/http"
    "net/rpc"
)

func main() {

    fmt.Println("Starting player")


	// create a hand as just an empty deck
    hand := new(playingcards.Deck)
    //tc := new(playingcards.TestCard)
    //thisCard := new(playingcards.Card)
    //hand.AddCard(*thisCard)
    //hand.Show()

    rpc.Register(hand)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp",":1234")
    if err != nil {
        fmt.Println("Error listening...")
    }
    go http.Serve(l,nil)

    fmt.Println("Waiting for hand to be dealt")
    for {
    }
}
