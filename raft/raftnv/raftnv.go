package raftnv

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type RaftLogEntry struct {
	Term  int
	Value string
}

type NVState struct {
	JSONFilename string         `json:"-"`
	lock         sync.RWMutex   `json:"-"`
	CurrentTerm  int            `json:"term"`
	VotedFor     string         `json:"votedFor"`
	Log          []RaftLogEntry `json:"log"`
}

// setter function for term
func (nvs *NVState) SetTerm(p int) {
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.CurrentTerm = p
	nvs.WriteNVState()
}

// setter function for votedFor
func (nvs *NVState) SetVotedFor(p string) {
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.VotedFor = p
	nvs.WriteNVState()
}

func (nvs *NVState) ReadNVState() error {

	// load the json file if it exists, otherwise load default initial state values
	filestat, err := os.Stat(nvs.JSONFilename)
	if err == nil {

		// JSON FILE EXISTS

		// make sure this path isn't a directory
		if filestat.IsDir() {
			return errors.New("Error: " + nvs.JSONFilename + " is a directory!")
		}

		// load json file
		log.Println("File ", nvs.JSONFilename, " exists, loading nvstate from it...")
		infile, err := os.OpenFile(nvs.JSONFilename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer infile.Close()

		// decode json file
		decoder := json.NewDecoder(infile)
		err = decoder.Decode(nvs)
		if err != nil {
			return err
		}

	} else {

		// JSON FILE DOES NOT EXIST
		log.Println("File ", nvs.JSONFilename, " does not exist, so setting nvstate elements to default initial values...")

		// we don't really need to do anything here since the zero values are fine
		//st.Term = 0
		//st.LeaderID = ""
		//st.Log = append(st.Log, "first command", "second command")
	}

	// done
	return nil

}

// now create an object and write to a json file
func (nvs *NVState) WriteNVState() error {

	log.Println("Writing nvstate to file...")
	outfile, _ := os.OpenFile(nvs.JSONFilename, os.O_CREATE, os.ModePerm)
	defer outfile.Close()
	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	encoder.Encode(nvs)

	// done
	return nil
}

func (nvs *NVState) AppendLogEntry(term int, value string) {
	entry := RaftLogEntry{
		Term:  term,
		Value: value,
	}
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.Log = append(nvs.Log, entry)
	nvs.WriteNVState()
}

func (nvs *NVState) GetLogTerm(index int) int {
	if index < 1 || index > len(nvs.Log) {
		return -1
	}
	return nvs.Log[index-1].Term
}
