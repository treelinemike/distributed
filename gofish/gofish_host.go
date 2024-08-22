package main

import (
	"engg415/playingcards"
	"fmt"
    "net/rpc"
    "math/rand"
    "time"
)

func main() {
    
    var j int // this is silly - RPC interface requires reply var even if not used...
    var err error

    // pull time as a seed for the rng
    // TODO: check whether we can get consistent performance with a specified seed
    // if so could use this for grading purposes
    rand.Seed(time.Now().UnixNano()) //https://stackoverflow.com/questions/12321133/how-to-properly-seed-random-number-generator

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
            // TODO: take top card and put it back if RPC fails
            c := deck.TakeTopCard()
            fmt.Printf("Taking top card: " + c.String()+"\n")
            err = player.Call("Deck.AddCardRPC",deck.TakeTopCard(),&j)
            if err != nil {
                fmt.Println("Error dealing card: "+err.Error())
            }
        }
    }

    // show remaining deck
    fmt.Printf("All hands dealt. Remaining deck has %d cards:\n",deck.NumCards())
    deck.Show()
}
