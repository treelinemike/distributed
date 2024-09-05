package main

import (
	"engg415/gofish/gfcommon"

	"log"
	"net"
	"net/http"
	"net/rpc"
	"slices"
)

func main() {

	var j int // this is silly - net/rpc interface requires reply var even if not used...
	var ret gfcommon.GFPlayerReturn
	var err error

	// format log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// load host and player IP addresses and ports from YAML config file
	var HostIP gfcommon.NetworkAddress
	var PlayerIP = make([]gfcommon.NetworkAddress, 0)
	var ConnectedPlayerIP = make([]gfcommon.NetworkAddress, 0)
	handSizes := make([]int, 0)
	numBooks := make([]int, 0)
	err = gfcommon.LoadGFGameConfig("gofish_config.yaml", &HostIP, &PlayerIP)
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
	log.Println("Ready to play")

	// create and shuffle a standard deck
	deck.Create()
	log.Printf("Deck created with %d cards: %s\n", deck.NumCards(), deck.String())
	deck.Shuffle()
	log.Printf("Deck shuffled: %s", deck.String())

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
	// TODO: do we want two at minimum?
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

		// set number of books and cards in hand to zero
		numBooks = append(numBooks, 0)
		handSizes = append(handSizes, 0)

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

	// deal cards to each player off the deck
	var numCardsToDeal int
	if len(players) < 3 {
		numCardsToDeal = 7
	} else {
		numCardsToDeal = 5
	}
	log.Println("Dealing cards...")
	for i := 0; i < numCardsToDeal; i++ {
		for playerIdx, player := range players {

			// take the top card off the deck and try to deal it to a player
			c := deck.TakeTopCard()
			err = player.Call("GFPlayerAPI.AddCardToHand", c, &j)

			// if dealing fails put the card back on the TOP of the deck
			if err != nil {
				log.Println("Error dealing card: " + err.Error())
				deck.AddCard(c)
				log.Printf("Card returned to deck, which now has %d cards: %s\n", deck.NumCards(), deck.String())

			} else {
				handSizes[playerIdx] += 1
			}
		}
	}

	// ask each player to print the hand it was dealt (to log)
	for playerIdx, player := range players {
		err = player.Call("GFPlayerAPI.PrintHand", j, &j)
		if err != nil {
			log.Printf("Could not get %s to print its hand\n", ConnectedPlayerIP[playerIdx].Address)
		}
	}

	// show remaining deck
	log.Printf("All hands dealt. Remaining deck has %d cards: %s\n", deck.NumCards(), deck.String())

	// cycle through each player
	//doneflag := false
	doneflag := false
	playerIdx := 0
	for !doneflag {
		log.Printf("Activiating player %d (%s)\n", playerIdx, ConnectedPlayerIP[playerIdx].Address)

		// interestingly struct fields with zero value aren't included in gob encoding
		// so we need to reset return struct fields to zero before calling an RPC
		// this was a terrible debug, but is actually confirmed here: https://github.com/golang/go/issues/8997
		ret.NumBooks = 0
		ret.NumCardsInHand = 0
		err = players[playerIdx].Call("GFPlayerAPI.TakeTurn", j, &ret)
		if err != nil {
			log.Print(err)
			log.Fatalf("Could not exectue TakeTurn RPC for player %d (%s)\n", playerIdx, ConnectedPlayerIP[playerIdx].Address)
		}
		log.Printf("Turn complete for player %d (%s): has %d books and %d cards in hand\n", playerIdx, ConnectedPlayerIP[playerIdx].Address, ret.NumBooks, ret.NumCardsInHand)
		numBooks[playerIdx] = ret.NumBooks
		handSizes[playerIdx] = ret.NumCardsInHand

		if (deck.NumCards() == 0) && (slices.Max(handSizes) == 0) {
			doneflag = true
		}

		// move on to next player
		playerIdx += 1
		if playerIdx >= len(players) {
			playerIdx = 0
		}
	}

	// figure out who won
	var winStatus int
	winString := "Game over! Won by: "
	winningNumBooks := slices.Max(numBooks)
	numWinners := 0
	for _, nb := range numBooks {
		if nb == winningNumBooks {
			numWinners++
		}
	}
	for playerIdx, player := range players {
		winStatus = 0
		sep := ""
		if numBooks[playerIdx] == winningNumBooks {
			winStatus = numWinners
			winString += sep + ConnectedPlayerIP[playerIdx].Address
			sep = " and "
		}
		err = player.Call("GFPlayerAPI.EndGame", winStatus, &j)
		if err != nil {
			log.Fatalf("Could not exectue EndGame RPC for player %d (%s)\n", playerIdx, ConnectedPlayerIP[playerIdx].Address)
		}
	}
	log.Println(winString)

}
