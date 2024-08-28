package main

import (
	"engg415/playingcards"
	//"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func main() {

	var j int // this is silly - RPC interface requires reply var even if not used...
	var err error

	// TODO: CHECK FOR IP ADDRESSES
	ifaces, err := net.Interfaces()
	fmt.Println("error %s",err)
    // handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("Error\n")
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// process IP address
				log.Printf("IP addr: %s", ip)
			}
		}

		// create and shuffle a standard deck
		deck := new(playingcards.Deck)
		deck.Create()
		log.Printf("Deck created with %d cards:\n", deck.NumCards())
		deck.Show()
		deck.Shuffle()
		log.Println("Deck shuffled:")
		deck.Show()

		// list of all players
		port := "1234"
		players_ip := [...]string{"192.168.10.10", "192.168.10.20"}

		// open network connections to all players
		log.Println("Connecting to remote player servers...")
		var players []*rpc.Client
		for _, ip := range players_ip {
			client, err := rpc.DialHTTP("tcp", ip+":"+port)
			if err != nil {
				log.Println("Error connecting to remote server: " + ip)
				continue
			}
			players = append(players, client)
			err = client.Call("Deck.ResetDeckRPC", j, &j)
		}

		// deal seven cards to each player
		log.Println("Dealing cards...")
		for i := 0; i < 7; i++ {
			for _, player := range players {

				// take the top card off the deck and try to deal it to a player
				c := deck.TakeTopCard()
				err = player.Call("Deck.AddCardRPC", c, &j)

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
}
