package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type nvstate struct {
	CurrentTerm int      `json:"term"`
	VotedFor    string   `json:"votedFor"`
	Log         []string `json:"log"`
}

// globals for non-volatile state
var jsonfilename string
var st nvstate
var stlock sync.RWMutex // Go doesn't like including lock in state struct due to json marshalling

// setter function for term
func setterm(p int) {
	stlock.Lock()
	defer stlock.Unlock()
	st.CurrentTerm = p
	writenvstate()
}

// setter function for leaderID
func setleaderid(p string) {
	stlock.Lock()
	defer stlock.Unlock()
	st.VotedFor = p
	writenvstate()
}

func readnvstate() error {

	// load the json file if it exists, otherwise load default initial state values
	filestat, err := os.Stat(jsonfilename)
	if err == nil {

		// JSON FILE EXISTS

		// make sure this path isn't a directory
		if filestat.IsDir() {
			return errors.New("Error: " + jsonfilename + " is a directory!")
		}

		// load json file
		log.Println("File ", jsonfilename, " exists, loading nvstate from it...")
		infile, err := os.OpenFile(jsonfilename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer infile.Close()

		// decode json file
		decoder := json.NewDecoder(infile)
		err = decoder.Decode(&st)
		if err != nil {
			return err
		}

	} else {

		// JSON FILE DOES NOT EXIST
		log.Println("File ", jsonfilename, " does not exist, so setting nvstate elements to default initial values...")

		// we don't really need to do anything here since the zero values are fine
		//st.Term = 0
		//st.LeaderID = ""
		//st.Log = append(st.Log, "first command", "second command")
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
