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
		switch e.Type {
		case 0: // wall types
			for _, v := range e.Data {
				var wt mazeviz.Walltype
				switch v {
				case 0:
					wt = mazeviz.W_none
				case 1:
					wt = mazeviz.W_latent
				case 2:
					wt = mazeviz.W_observed
				case 3:
					wt = mazeviz.W_phantom
				}
				p.Walltypes = append(p.Walltypes, wt)
			}
		case 100: // cell types
			for _, v := range e.Data {
				var ct mazeviz.Celltype
				switch v {
				case 0:
					ct = mazeviz.C_none
				case 1:
					ct = mazeviz.C_goal
				case 2:
					ct = mazeviz.C_start
				}
				p.Celltypes = append(p.Celltypes, ct)
			}
		case 101: // cell values
			for _, v := range e.Data {
				p.Cellvals = append(p.Cellvals, float32(v))
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
