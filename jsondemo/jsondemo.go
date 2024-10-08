package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type nvstate struct {
	Term     int      `json:"term"`
	LeaderID string   `json:"leaderID"`
	Log      []string `json:"log"`
}

var st nvstate

func main() {
	st.Term = 12
	st.LeaderID = "Server1"
	st.Log = append(st.Log, "first command", "second command")

	outfile, _ := os.OpenFile("nvstate.json", os.O_CREATE, os.ModePerm)
	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	encoder.Encode(st)
	outfile.Close()

	infile, _ := os.OpenFile("nvstate.json", os.O_RDONLY, os.ModePerm)
	decoder := json.NewDecoder(infile)
	var readnv nvstate
	decoder.Decode(&readnv)
	infile.Close()

	fmt.Println("Read from json: ", readnv)

}
