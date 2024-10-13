package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type nvstate struct {
	Term     int      `json:"term"`
	LeaderID string   `json:"leaderID"`
	Log      []string `json:"log"`
	Lock     sync.RWMutex
}

var jsonfilename string
var st nvstate

func setterm(p int) {
	st.Lock.Lock()
	defer st.Lock.Unlock()
	st.Term = p
	writenvstate()
}

func setleaderid(p string) {
	st.Lock.Lock()
	defer st.Lock.Unlock()
	st.LeaderID = p
	writenvstate()
}

func readnvstate() error {

	// load the json file if it exists, otherwise load default initial state values
	filestat, err := os.Stat(jsonfilename)
	if err == nil {

		// JSON FILE EXISTS

		// make sure this path isn't a directory
		if filestat.IsDir() {
			log.Fatal("Error: ", jsonfilename, " is a directory!")
		}

		// load json file into nvstate
		log.Println("File ", jsonfilename, " exists, loading nvstate from it...")
		infile, _ := os.OpenFile(jsonfilename, os.O_RDONLY, os.ModePerm)
		defer infile.Close()
		decoder := json.NewDecoder(infile)
		decoder.Decode(&st)

	} else {

		// JSON FILE DOES NOT EXIST
		log.Println("File ", jsonfilename, " does not exist, so set nvstate to default initial state...")
		st.Term = 12
		st.LeaderID = "Server1"
		st.Log = append(st.Log, "first command", "second command")
	}

	// done
	return nil

}

// now create an object and write to a json file
func writenvstate() error {

	log.Println("Writing nvstate to file...")
	outfile, _ := os.OpenFile(jsonfilename, os.O_CREATE, os.ModePerm)
	defer outfile.Close()
	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	encoder.Encode(st)

	// done
	return nil
}
