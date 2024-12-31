package main

import (
	"log"
)

func main() {

	// load a maze configuration from json
	readmaze, err := readjsonmaze("mazein.json")
	if err != nil {
		log.Fatalf("readjsonmaze error: %v\n", err)
	}
	log.Printf("Read: %v\n", readmaze)

	// generate a maze configuration
	writemaze := new(Mazedata)
	writemaze.Title = "test output"
	writemaze.Author = "kokko"
	writemaze.Description = "None"
	writemaze.M = 3
	writemaze.N = 3
	newelement := new(Mazeelement)
	newelement.Type = 100
	newelement.Description = "No description"
	newelement.Data = []float64{0, 1, 0, 0, 1, 0, 0, 1, 0}
	writemaze.Elements = append(writemaze.Elements, *newelement)

	// write the maze configuration to json
	err = writejsonmaze("mazeout.json", *writemaze)
	if err != nil {
		log.Fatalf("writejsonmaze error: %v\n", err)
	}
	log.Printf("Wrote: %v\n", *writemaze)

}
