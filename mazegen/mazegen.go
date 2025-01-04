package main

import (
	"engg415/mazeviz"
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// command line arguments
	mpflag := flag.Int("m", 0, "number of rows in the maze")
	npflag := flag.Int("n", 0, "number of columnss in the maze")
	ipflag := flag.String("i", "", "name of input json file")
	opflag := flag.String("o", "", "name of output json file")
	flag.Parse()

	// catch default case(s) with no args
	// e.g. exectued wtih go run .
	if *ipflag == "" {
		switch {
		case *mpflag == 0:
			*mpflag = 16
		case *npflag == 0:
			*npflag = 16
		}
	}
	if *opflag == "" {
		*opflag = "mazeout.json"
	}

	// create a new game
	game, err := mazeviz.Newgame()
	if err != nil {
		log.Fatal(err)
	}

	// load maze from json
	var ww, wh int
	if *ipflag != "" {
		ww, wh, err = game.Loadmaze(*ipflag)
	} else {
		ww, wh, err = game.Newmaze(*mpflag, *npflag)
	}
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
	err = game.Savemaze(*opflag)
	if err != nil {
		log.Fatalf("error saving maze: %v\n", err)
	}
}
