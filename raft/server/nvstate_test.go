package main

import (
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
	jsonfilename = "notafile.json"

	// confirm that file does not exist
	_, err := os.Stat(jsonfilename)
	if err == nil {
		t.Fatalf("File %s exists! Delete it and rerun test\n", jsonfilename)
	}
	t.Logf("Confirmed that file %s does not exist\n", jsonfilename)

	// attempt to read non-existent file
	err = readnvstate()
	if err != nil {
		t.Fatalf("readnvstate() failed with error: %v\n", err)
	}

	// make sure that default values are set
	if st.Term != 0 {
		t.Fatalf("Expected default Term of 0, got %d\n", st.Term)
	}
	if st.LeaderID != "" {
		t.Fatalf("Expected default LeaderID of empty string, got %s\n", st.LeaderID)
	}
	if st.Log != nil {
		t.Fatalf("Expected default Log of empty slice, got %v\n", st.Log)
	}
	t.Logf("Default values set correctly: Term=%d, LeaderID='%s', Log = nil slice\n", st.Term, st.LeaderID)

}

// NVSTATE: attempt to write and read nvstate to/from file
func TestWriteAndReadNVState(t *testing.T) {
	jsonfilename = "testnvstate.json"
	defer os.Remove(jsonfilename) // clean up file when done
	// set some values
	setterm(42)
	setleaderid("TestServer")
	st.Log = append(st.Log, "command1", "command2", "command3")
	t.Logf("Set nvstate to: Term=%d, LeaderID='%s', Log=%v\n", st.Term, st.LeaderID, st.Log)

	// write to file
	err := writenvstate()
	if err != nil {
		t.Fatalf("writenvstate() failed with error: %v\n", err)
	}

	// reset in-memory nvstate
	st.Term = 0
	st.LeaderID = ""
	st.Log = nil
	t.Logf("Reset in-memory nvstate to: Term=%d, LeaderID='%s', Log=%v\n", st.Term, st.LeaderID, st.Log)

	// read from file
	err = readnvstate()
	if err != nil {
		t.Fatalf("readnvstate() failed with error: %v\n", err)
	}

	// verify values
	if st.Term != 42 {
		t.Fatalf("Expected Term of 42, got %d\n", st.Term)
	}
	if st.LeaderID != "TestServer" {
		t.Fatalf("Expected LeaderID of 'TestServer', got '%s'\n", st.LeaderID)
	}
	if len(st.Log) != 3 {
		t.Fatalf("Expected Log of length 3, got %d\n", len(st.Log))
	}
	for i, cmd := range []string{"command1", "command2", "command3"} {
		if st.Log[i] != cmd {
			t.Fatalf("Expected Log[%d] to be '%s', got '%s'\n", i, cmd, st.Log[i])
		}
	}
	t.Logf("Values verified correctly: Term=%d, LeaderID='%s', Log=%v\n", st.Term, st.LeaderID, st.Log)
}
