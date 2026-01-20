package raftnv

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// use TestMain() to set logging to microseconds for better resolution
func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	os.Exit(m.Run())
}

// NVSTATE: attempt to read a non-existent file, this should succeed with default values set for nvstate
func TestLoadNonExistantFile(t *testing.T) {

	// create a state object
	var st NVState
	st.JSONFilename = "notafile.json"

	// confirm that file does not exist
	_, err := os.Stat(st.JSONFilename)
	if err == nil {
		t.Fatalf("File %s exists! Delete it and rerun test\n", st.JSONFilename)
	}
	t.Logf("Confirmed that file %s does not exist\n", st.JSONFilename)

	// attempt to read non-existent file
	err = st.ReadNVState()
	if err != nil {
		t.Fatalf("ReadNVState() failed with error: %v\n", err)
	}

	// make sure that default values are set
	if st.CurrentTerm != 0 {
		t.Fatalf("Expected default Term of 0, got %d\n", st.CurrentTerm)
	}
	if st.VotedFor != "" {
		t.Fatalf("Expected default VotedFor value of empty string, got %s\n", st.VotedFor)
	}
	if st.Log != nil {
		t.Fatalf("Expected default Log of empty slice, got %v\n", st.Log)
	}
	t.Logf("Default values set correctly: Term=%d, VotedFor='%s', Log = nil slice\n", st.CurrentTerm, st.VotedFor)

}

// NVSTATE: attempt to write and read nvstate to/from file
func TestWriteAndReadNVState(t *testing.T) {

	// create a state object
	var st NVState
	st.JSONFilename = "testnvstate.json"
	defer os.Remove(st.JSONFilename)

	// set some values
	st.SetTerm(42)
	st.SetVotedFor("TestServer")
	for i := 0; i <= 2; i++ {
		st.AppendLogEntry((42 + i), fmt.Sprintf("command%02d", i+1))
	}
	t.Logf("Set nvstate to: Term=%d, VotedFor='%s', Log=%v\n", st.CurrentTerm, st.VotedFor, st.Log)

	// write to file
	err := st.WriteNVState()
	if err != nil {
		t.Fatalf("WriteNVState() failed with error: %v\n", err)
	}

	// reset in-memory nvstate
	st.CurrentTerm = 0
	st.VotedFor = ""
	st.Log = nil
	t.Logf("Reset in-memory nvstate to: Term=%d, VotedFor='%s', Log=%v\n", st.CurrentTerm, st.VotedFor, st.Log)

	// read from file
	err = st.ReadNVState()
	if err != nil {
		t.Fatalf("ReadNVState() failed with error: %v\n", err)
	}

	// verify values
	if st.CurrentTerm != 42 {
		t.Fatalf("Expected Term of 42, got %d\n", st.CurrentTerm)
	}
	if st.VotedFor != "TestServer" {
		t.Fatalf("Expected VotedFor of 'TestServer', got '%s'\n", st.VotedFor)
	}
	if len(st.Log) != 3 {
		t.Fatalf("Expected Log of length 3, got %d\n", len(st.Log))
	}
	for i, cmd := range []string{"command01", "command02", "command03"} {
		if st.Log[i].Value != cmd {
			t.Fatalf("Expected Log[%d] to be '(%d,%s)', got '(%d,%s)'\n", i, i+42, cmd, st.Log[i].Term, st.Log[i].Value)
		}
	}
	t.Logf("Values verified correctly: Term=%d, VotedFor='%s', Log=%v\n", st.CurrentTerm, st.VotedFor, st.Log)
}
