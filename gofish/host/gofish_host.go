package main

import (
	"engg415/playingcards"
	//"errors"
	"fmt"
	"log"
	"net/rpc"
)

func main() {

	var j int // this is silly - RPC interface requires reply var even if not used...
	var err error

	// load host and player IP addresses and ports from YAML config file
	var HostIP NetworkAddress
	var PlayerIP = make([]NetworkAddress, 0)
	err = LoadGFConfig("gofish_config.yaml", &HostIP, &PlayerIP)
	if err != nil {
		log.Fatal(err)
	}

	// create and shuffle a standard deck
	deck := new(playingcards.Deck)
	deck.Create()
	log.Printf("Deck created with %d cards:\n", deck.NumCards())
	deck.Show()
	deck.Shuffle()
	log.Println("Deck shuffled:")
	deck.Show()

	// open network connections to all players
	log.Println("Connecting to remote player servers...")
	var players []*rpc.Client
	for _, playerip := range PlayerIP {
		client, err := rpc.DialHTTP("tcp", playerip.Address+":"+playerip.Port)
		if err != nil {
			log.Println("Error connecting to remote server: " + playerip.Address)
			continue
		}
		players = append(players, client)
		err = client.Call("GFPlayerAPI.ResetHand", j, &j)
		if err != nil {
			fmt.Println("Couldn't reset deck")
		}
	}

	// deal seven cards to each player
	log.Println("Dealing cards...")
	for i := 0; i < 7; i++ {
		for _, player := range players {

			// take the top card off the deck and try to deal it to a player
			c := deck.TakeTopCard()
			err = player.Call("GFPlayerAPI.AddCardToHand", c, &j)

			// if dealing fails put the card back on the TOP of the deck
			if err != nil {
				fmt.Println("Error dealing card: " + err.Error())
				deck.AddCard(c)
				log.Printf("Card returned to deck, which now has %d cards:\n", deck.NumCards())
				deck.Show()
			}
		}
	}

	// show remaining deck
	log.Printf("All hands dealt. Remaining deck has %d cards:\n", deck.NumCards())
	deck.Show()
}
