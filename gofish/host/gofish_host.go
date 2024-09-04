package main

import (
	"engg415/gofish/gfcommon"

	//"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {

	var j int // this is silly - net/rpc interface requires reply var even if not used...
	var err error

	// load host and player IP addresses and ports from YAML config file
	var HostIP gfcommon.NetworkAddress
	var PlayerIP = make([]gfcommon.NetworkAddress, 0)
	var ConnectedPlayerIP = make([]gfcommon.NetworkAddress, 0)
	err = LoadGFConfig("gofish_config.yaml", &HostIP, &PlayerIP)
	if err != nil {
		log.Fatal(err)
	}

	// create and register host API
	log.Println("Registering host API")
	thisapi := new(GFHostAPI)
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+HostIP.Port)
	if err != nil {
		log.Fatal("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)
	fmt.Println("Ready to play")

	// create and shuffle a standard deck
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

		// connect to player
		client, err := rpc.DialHTTP("tcp", playerip.Address+":"+playerip.Port)
		if err != nil {
			log.Println("Error connecting to remote server: " + playerip.Address)
			continue
		}

		// note: players and ConnectedPlayerIP are always 1:1
		players = append(players, client)
		ConnectedPlayerIP = append(ConnectedPlayerIP, playerip)
	}

	// make sure we have some players
	if len(players) == 0 {
		log.Fatal("No players connected")
	}

	// configure each player
	for i, player := range players {

		log.Printf("Setting config for player: %s\n", ConnectedPlayerIP[i].Address)

		// reset player hand
		err = player.Call("GFPlayerAPI.ResetHand", j, &j)
		if err != nil {
			log.Fatal("Couldn't reset deck")
		}

		// assemble list of other players
		playerConfig := new(gfcommon.GFPlayerConfig)
		playerConfig.Host = HostIP
		for j, otherplayerip := range ConnectedPlayerIP {
			if j != i {
				playerConfig.OtherPlayers = append(playerConfig.OtherPlayers, otherplayerip)
			}
		}

		// set list of other players via RPC
		err = player.Call("GFPlayerAPI.SetConfig", playerConfig, &j)
		if err != nil {
			log.Fatal("Couldn't set player config")
		}
	}

	// deal seven cards to each player off the deck
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

	// cycle through each player
	//doneflag := false
	doneflag := false
	playerIdx := 0
	rounds := 0
	for !doneflag {
		fmt.Printf("Activiating player %d (%s)\n", playerIdx, ConnectedPlayerIP[playerIdx].Address)
		err = players[playerIdx].Call("GFPlayerAPI.TakeTurn", j, &j)
		if err != nil {
			log.Print(err)
			log.Fatalf("Could not exectue TakeTurn RPC for player %d (%s)\n", playerIdx, ConnectedPlayerIP[playerIdx].Address)
		}

		playerIdx++
		if playerIdx == len(players) {
			playerIdx = 0
			rounds++
		}

		if rounds == 4 {
			doneflag = true
		}
	}
}
