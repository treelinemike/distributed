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

	// catch default case(s) with no args
	// or when exectued wtih go run .
	if *ipflag == "" {
		log.Fatal("Must specify maze json file")
	}

	// create a new game
	game, err := mazeviz.Newgame()
	if err != nil {
		log.Fatal(err)
	}

	// load maze from json or specify default size
	ww, wh, err = game.Loadmaze(*ipflag)
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

	// save maze in json file
	// when game window is closed
	err = game.Savemaze(*opflag)
	if err != nil {
		log.Fatalf("error saving maze: %v\n", err)
	}
}
