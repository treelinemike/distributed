package main

import (
	"engg415/mazeviz"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// create a new game
	game, err := mazeviz.Newgame()
	if err != nil {
		log.Fatal(err)
	}

	// load maze from json
	ww, wh, err := game.Loadmaze("mazein_1.json")
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
	err = game.Savemaze("mazeout.json")
	if err != nil {
		log.Fatalf("error saving maze: %v\n", err)
	}

}
