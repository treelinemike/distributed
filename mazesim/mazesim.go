package main

import (
	"engg415/mazeviz"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// load a maze configuration from json
	readmaze, err := readjsonmaze("mazein.json")
	if err != nil {
		log.Fatalf("readjsonmaze error: %v\n", err)
	}
	log.Printf("Read: %v\n", readmaze)

	//mazeviz.Setparams(int(readmaze.M), int(readmaze.N))

	p := new(mazeviz.Params)

	p.Setparams(int(readmaze.M), int(readmaze.N))
	for _, e := range readmaze.Elements {
		if e.Type == 0 {
			for _, v := range e.Data {
				p.Walls = append(p.Walls, int32(v))
			}
		}
		if e.Type == 100 {
			for _, v := range e.Data {
				p.Cellvals = append(p.Walls, int32(v))
			}
		}

	}

	//p.Setparams(16, 16)

	game, err := mazeviz.NewGame(*p)
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(p.WW, p.WH)
	ebiten.SetWindowTitle("Maze Generator")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	// generate a maze configuration
	writemaze := new(Mazedata)
	writemaze.Title = "test output"
	writemaze.Author = "kokko"
	writemaze.Description = ""
	writemaze.M = 3
	writemaze.N = 3
	newelement := new(Mazeelement)
	newelement.Type = 100
	newelement.Description = ""
	newelement.Data = []float32{0, 1, 0, 0, 1, 0, 0, 1, 0}
	writemaze.Elements = append(writemaze.Elements, *newelement)

	// write the maze configuration to json
	err = writejsonmaze("mazeout.json", *writemaze)
	if err != nil {
		log.Fatalf("writejsonmaze error: %v\n", err)
	}
	log.Printf("Wrote: %v\n", *writemaze)

}
