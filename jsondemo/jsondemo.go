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
	if err == nil {
		// make sure this path isn't a directory
		if filestat.IsDir() {
			log.Fatal("Error: ", jsonfilename, " is a directory!")
		}

		// load json file into nv struct
		log.Println("File ", jsonfilename, " exists, loading it...")
		infile, _ := os.OpenFile(jsonfilename, os.O_RDONLY, os.ModePerm)
		decoder := json.NewDecoder(infile)
		decoder.Decode(&st)
		infile.Close()

	} else {
		log.Println("File ", jsonfilename, " does not exist, so initialize nv struct...")
		st.Term = 12
		st.LeaderID = "Server1"
		st.Log = append(st.Log, "first command", "second command")
	}

	// show nv struct
	fmt.Println("NV struct contents: ", st)

	// now create an object and write to a json file
	log.Println("Writing nv struct to file...")
	outfile, _ := os.OpenFile(jsonfilename, os.O_CREATE, os.ModePerm)
	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	encoder.Encode(st)
	outfile.Close()

}
