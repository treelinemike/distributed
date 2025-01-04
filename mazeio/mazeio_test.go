package mazeio

import (
	"testing"
)

func TestRead(t *testing.T) {
	// load a maze configuration from json
	_, err := Readjsonmaze("mazeio_test.json")
	if err != nil {
		t.Fatalf("readjsonmaze error: %v\n", err)
	}
}

func TestWrite(t *testing.T) {

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
	err := Writejsonmaze("mazeout.json", *writemaze)
	if err != nil {
		t.Fatalf("writejsonmaze error: %v\n", err)
	}
}
