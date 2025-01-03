package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

// write maze struct to json
func writejsonmaze(filename string, writemaze Mazedata) error {

	err := verifymazestructs(writemaze)
	if err != nil {
		return err
	}
	outjson, err := json.Marshal(&writemaze) //json.MarshalIndent(writemaze, "", "    ")
	if err != nil {
		return err
	}
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = outfile.Write(outjson)
	if err != nil {
		return err
	}
	err = outfile.Close()
	if err != nil {
		return err
	}

	return nil
}

// read maze struct from json
func readjsonmaze(filename string) (Mazedata, error) {
	fid, err := os.Open(filename)
	if err != nil {
		return Mazedata{}, err
	}
	jsondata, err := io.ReadAll(fid)
	if err != nil {
		return Mazedata{}, err
	}
	fid.Close()

	readmaze := new(Mazedata)
	json.Unmarshal(jsondata, readmaze)
	err = verifymazestructs(*readmaze)
	if err != nil {
		return Mazedata{}, err
	}
	return *readmaze, nil
}

// check that we have the correct number of elements in all data slices
func verifymazestructs(maze Mazedata) error {

	for _, e := range maze.Elements {

		var expectedlen int32
		if e.Type < 100 {
			expectedlen = (2*maze.M+1)*maze.N + maze.M
		} else {
			expectedlen = maze.M * maze.N
		}

		if len(e.Data) != int(expectedlen) {
			log.Printf("[type %d] incorrect data length: %d\n", e.Type, len(e.Data))
			return errors.New("data length inconsistent with maze size and type")
		}
	}

	return nil
}

// nested json structure for maze state storage
type Mazedata struct {
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	Description string        `json:"description"`
	M           int32         `json:"m"`
	N           int32         `json:"n"`
	Elements    []Mazeelement `json:"elements"`
}

// maze element storage
type Mazeelement struct {
	Type        int32     `json:"type"`
	Description string    `json:"description"`
	Data        []float32 `json:"data"`
}
