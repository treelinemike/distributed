package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type nvstate struct {
	Term     int      `json:"term"`
	LeaderID string   `json:"leaderID"`
	Log      []string `json:"log"`
}

var st nvstate

func main() {

	// name of the json file to which we're storing our nonvolatile data
	jsonfilename := "nvstate.json"

	// ensure that we're starting from scratch with no json file in local directory
	filestat, err := os.Stat(jsonfilename)
	if err != nil {
		log.Println("File ", jsonfilename, " does not exist...")
	} else {
		if filestat.IsDir() {
			log.Fatal("Error: ", jsonfilename, " is a directory!")
		}
		log.Println("File ", jsonfilename, " alread exists, deleting it...")
		err := os.Remove(jsonfilename)
		if err != nil {
			log.Fatal("Could not remove file")
		}
	}

	// now create an object and write to a json file
	st.Term = 12
	st.LeaderID = "Server1"
	st.Log = append(st.Log, "first command", "second command")
	outfile, _ := os.OpenFile(jsonfilename, os.O_CREATE, os.ModePerm)
	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	encoder.Encode(st)
	outfile.Close()

	infile, _ := os.OpenFile(jsonfilename, os.O_RDONLY, os.ModePerm)
	decoder := json.NewDecoder(infile)
	var readnv nvstate
	decoder.Decode(&readnv)
	infile.Close()

	fmt.Println("Read from json: ", readnv)

}
