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

// Note: CurrentTerm, VotedFor, and Log must be exported to be encoded/decoded by json package,
// but we want to be careful about concurrent access, so we use a RWMutex to protect them which requires
// us to be good about using getter and setter functions for access
// TODO: use a separate data structure for json encoding/decoding to avoid exporting these fields?
type NVState struct {
	jsonFilename string         `json:"-"`
	lock         sync.RWMutex   `json:"-"`
	CurrentTerm  int            `json:"term"`
	VotedFor     string         `json:"votedFor"`
	Log          []RaftLogEntry `json:"log"`
}

// get log length
func (nvs *NVState) LogLength() int {
	nvs.lock.RLock()
	defer nvs.lock.RUnlock()
	return len(nvs.Log)
}

// jsonFilename getter function
func (nvs *NVState) GetJSONFilename() string {
	nvs.lock.RLock()
	defer nvs.lock.RUnlock()
	return nvs.jsonFilename
}

// jsonFilename setter function
func (nvs *NVState) SetJSONFilename(filename string) {
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.jsonFilename = filename
}

// CurrentTerm getter function
func (nvs *NVState) GetCurrentTerm() int { // not idomatic Go, but we need to export both field and method
	nvs.lock.RLock()
	defer nvs.lock.RUnlock()
	return nvs.CurrentTerm
}

// CurrentTerm setter function
func (nvs *NVState) SetCurrentTerm(p int) {
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.CurrentTerm = p
	nvs.WriteNVState()
}

// VotedFor getter function
func (nvs *NVState) GetVotedFor() string {
	nvs.lock.RLock()
	defer nvs.lock.RUnlock()
	return nvs.VotedFor
}

// VotedFor setter function
func (nvs *NVState) SetVotedFor(p string) {
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.VotedFor = p
	nvs.WriteNVState()
}

// load state from non-volatile storage (json file)
func (nvs *NVState) ReadNVState() error {

	// load the json file if it exists, otherwise load default initial state values
	filestat, err := os.Stat(nvs.jsonFilename)
	if err == nil {

		// JSON FILE EXISTS

		// make sure this path isn't a directory
		if filestat.IsDir() {
			return errors.New("Error: " + nvs.jsonFilename + " is a directory!")
		}

		// load json file
		log.Println("File ", nvs.jsonFilename, " exists, loading nvstate from it...")
		infile, err := os.OpenFile(nvs.jsonFilename, os.O_RDONLY, os.ModePerm)
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
		log.Println("File ", nvs.jsonFilename, " does not exist, so setting nvstate elements to default initial values...")

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
	outfile, err := os.OpenFile(nvs.jsonFilename, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer outfile.Close()
	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(nvs)
	if err != nil {
		return err
	}

	// done
	return nil
}

func (nvs *NVState) AppendLogEntry(term int, value string) error {
	entry := RaftLogEntry{
		Term:  term,
		Value: value,
	}
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.Log = append(nvs.Log, entry)
	err := nvs.WriteNVState()
	return err
}

func (nvs *NVState) GetLogEntry(index int) (RaftLogEntry, error) {
	if index < 1 || index > len(nvs.Log) {
		return RaftLogEntry{}, errors.New("Index out of bounds")
	}
	return nvs.Log[index-1], nil
}

func (nvs *NVState) GetLogEntryTerm(index int) (int, error) {
	if index < 1 || index > len(nvs.Log) {
		return 0, errors.New("Index out of bounds")
	}
	return nvs.Log[index-1].Term, nil
}

func (nvs *NVState) PruneLog(index int) error {
	if index < 1 || index > len(nvs.Log) {
		return errors.New("Index out of bounds")
	}
	nvs.lock.Lock()
	defer nvs.lock.Unlock()
	nvs.Log = nvs.Log[:index]
	err := nvs.WriteNVState()
	return err
}
