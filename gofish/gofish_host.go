package main

import (
	"engg415/playingcards"
	"fmt"
    "net/rpc"
)

func main() {

	// create and shuffle a standard deck
	deck := new(playingcards.Deck)
	deck.Create()
	deck.Create()
    deck.Shuffle()
	fmt.Println("Shuffled deck:")
	deck.Show()
	fmt.Println("Top card is: ", deck.TakeTopCard().String())
	fmt.Println("Deck is now:")
	deck.Show()
	deck.Shuffle()
	fmt.Println("After shuffling again:")
	deck.Show()

    // now we deal one card to test RPC
    client, err := rpc.DialHTTP("tcp","192.168.10.10:1234")
    if err != nil {
        fmt.Println("Error connecting to remote server...")
    }
    /*a := new(playingcards.MyArgType)
    a.X = 1
    a.ThisSuit = playingcards.Hearts
    //import "errors"
    j := 2
    err = client.Call("TestCard.TestRPC",a,&j)
    */
    card := new(playingcards.Card)
    var j int
    err = client.Call("Deck.AddCardRPC",card,&j)
    if err!= nil {
        fmt.Println("Error executing RPC: "+err.Error())
    }

}
