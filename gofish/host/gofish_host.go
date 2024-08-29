package main

import (
	"engg415/playingcards"
	//"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type (
	NetworkAddress struct {
		Address string `mapstructure:"Address"`
		Port    int    `mapstructure:"Port"`
	}
)

type GFConfig struct {
	Host    NetworkAddress            `mapstructure:"Host"`
	Players map[string]NetworkAddress `mapstructure:"Players"`
}

func main() {

	var j int // this is silly - RPC interface requires reply var even if not used...
	var err error

	configFile, err := os.ReadFile("gofish_config.yaml")
	if err != nil {
		log.Fatal("Cannot read config file!")
	}

	var cfg GFConfig
	var yamlIn interface{}

	err = yaml.Unmarshal(configFile, &yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &cfg})
	err = decoder.Decode(yamlIn)
	if err != nil {
		log.Fatal(err)
	}

	// Print out the new struct
	fmt.Printf("%+v\n", cfg)

	fmt.Printf("Config file host address: %s:%d\n", cfg.Host.Address, cfg.Host.Port)
	for player := range cfg.Players {

		fmt.Printf("Player name: %s on %s:%d\n", player, cfg.Players[player].Address, cfg.Players[player].Port)

	}
	fmt.Printf("first player address: %s\n", cfg.Players["Player3"].Address)

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
	players_ip := [...]string{"localhost", "192.168.10.10", "192.168.10.20"}

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
			err = player.Call("GFPlayerAPI.AddCard", c, &j)

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
