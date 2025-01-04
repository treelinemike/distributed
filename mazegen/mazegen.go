// mazegen.go
// App for generating and editing 2D maze confugrations
// Graphics handled by ebitenengine v2
// Uses custom mazeviz package, which also runs maze simulator
// Author:   M. Kokko
// Modified: 04-Jan-2024
//
// Usage examples:
// mazegen                    edit an empty 16x16 maze, saves as mazeout.json on exit
// mazegen -o mymaze.json     edit an empty 16x16 maze, saves as mymaze.json on exit
// mazegen -m=3 -n=5          edit an empty maze with 3 rows and 5 columns
// mazegen -i mazein_1.json   edit existing maze from mazein_1.json

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
	// or when exectued wtih go run .
	if *ipflag == "" {
		if *mpflag == 0 {
			*mpflag = 16
		}
		if *npflag == 0 {
			*npflag = 16
		}
	} else {
		if *npflag != 0 || *mpflag != 0 {
			log.Fatal("Cannot specify both input file and maze dimensions")
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

	// load maze from json or specify default size
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
	// when game window is closed
	err = game.Savemaze(*opflag)
	if err != nil {
		log.Fatalf("error saving maze: %v\n", err)
	}
}
