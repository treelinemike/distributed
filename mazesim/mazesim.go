// mazesim.go

package main

import (
	"engg415/mazeviz"
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// command line arguments
	ipflag := flag.String("i", "", "name of input json file")
	flag.Parse()
	if *ipflag == "" {
		*ipflag = "defaultmaze.json"
		log.Println("No input maze specified, usinsg default")
	}

	// create a new game
	game, err := mazeviz.Newgame()
	if err != nil {
		log.Fatal(err)
	}

	// load maze from json or specify default size
	ww, wh, err := game.Loadmaze(*ipflag)
	if err != nil {
		log.Fatalf("error loading maze: %v\n", err)
	}

	// configure window
	ebiten.SetWindowSize(ww, wh)
	ebiten.SetWindowTitle("Maze Generator")

	// run game
	if err = ebiten.RunGame(game); err != nil {
		log.Fatalf("error running game: %v\n", err)
	}

}
