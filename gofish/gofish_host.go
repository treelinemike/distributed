package main

import (
	"engg415/playingcards"
	"fmt"
    "net/rpc"
)

func main() {
    
    var j int // this is silly - RPC interface requires reply var even if not used...
    var err error

	// create and shuffle a standard deck
	deck := new(playingcards.Deck)
	deck.Create()
    fmt.Printf("Deck created with %d cards:\n",deck.NumCards())
    deck.Show()
    deck.Shuffle()
    fmt.Println("Deck shuffled:")
    deck.Show()

    // list of all players
    port := "1234"
    players_ip := [...]string{"192.168.10.10","192.168.10.20"}
    
    // open network connections to all players
    var players []*rpc.Client
    for _, ip := range players_ip {
        fmt.Println(ip)
        client, err := rpc.DialHTTP("tcp",ip+":"+port)
        if err != nil {
            fmt.Println("Error connecting to remote server: "+ip)
            continue
        }   
        players = append(players,client)
        err = client.Call("Deck.ResetDeckRPC",j,&j)
    }

    // deal seven cards to each player
    fmt.Println("Dealing cards...")
    for i:= 0; i<7; i++ {
        for _, player := range players {
            
            // take the top card off the deck and try to deal it to a player
            c := deck.TakeTopCard()
            fmt.Printf("Taking top card: " + c.String()+"\n")
            err = player.Call("Deck.AddCardRPC",c,&j)
            
            // if dealing fails put the card back on the TOP of the deck
            if err != nil {
                fmt.Println("Error dealing card: "+err.Error())
                deck.AddCard(c)
                fmt.Printf("Card returned to deck, which now has %d cards:\n",deck.NumCards())
                deck.Show()
            }
        }
    }

    // show remaining deck
    fmt.Printf("All hands dealt. Remaining deck has %d cards:\n",deck.NumCards())
    deck.Show()
}
