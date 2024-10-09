package main

import (
	"engg415/raft/common"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

type RaftAPI int

var electiontimer *time.Timer

// non-volitile storage --- how are we going to implement this?
var currentTerm int = 0 // can't do short declaration := at package level
var votedFor string = ""
var raftlog []string = make([]string, 0)

// TODO: struct members don't need to be exported if we're keeping in the server.go main package
type AEParams struct {
	term         int
	leaderId     int
	prevLogIndex int
	leaderCommit int
}

func (r *RaftAPI) AppendEntries(p AEParams, resp *int) error {
	fmt.Printf("AppendEntriesRPC: term %d\n", p.term)
	return nil
}

type RVParams struct {
	term         int
	candidateId  int
	lastLogIndex int
	lastLogTerm  int
}

type RVResp struct {
	term        int
	voteGranted bool
}

func (r *RaftAPI) RequestVote(p RVParams, resp *RVResp) error {
	fmt.Printf("RequestVoteRPC: term %d\n", p.term)
	return nil
}

func main() {

	var err error

	// format log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// load cluster configuration
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . configfile.yaml serverkey")
		return
	}
	filename := os.Args[1]
	selfkey := os.Args[2]
	servers := make(map[string]common.NetworkAddress)

	// load config
	log.Println("Loading Raft config file")
	var t common.Timeout
	common.LoadRaftConfig(filename, servers, &t, &jsonfilename)
	log.Printf("Config specifies election timeout range [%d, %d]\n", t.Min_ms, t.Max_ms)

	// make sure provided selfkey is in map from config file
	_, ok := servers[selfkey]
	if !ok {
		log.Fatal("Key ", selfkey, " is not in cluster config file")
	} else {
		log.Printf("Config specifies this server as %s:%s", servers[selfkey].Address, servers[selfkey].Port)
	}

	// load non-volatile state if it has been previously saved
	log.Println("Reading non-volatile state from file: ", jsonfilename)
	err = readnvstate()
	if err != nil {
		log.Fatal("Error initializing initial state")
	}

	// test setting term
	setterm(14)
	setleaderid("someserver")

	/*
		// test writing non-volatile state to JSON file
		log.Println("Writing non-volatile state to file: ", jsonfilename)
		err = writenvstate()
		if err != nil {
			log.Fatal("Error writing state")
		}
	*/

	// serve RaftAPI
	log.Println("Registering RPCs for access on port", servers[selfkey].Port)
	thisapi := new(RaftAPI)
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+servers[selfkey].Port)
	if err != nil {
		log.Fatal("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)

	// start a timer
	electiontimer = time.NewTimer(100 * time.Second)
	i := 0
	for {
		randtval := time.Duration(rand.IntN(t.Max_ms-t.Min_ms)+t.Min_ms) * time.Millisecond
		log.Printf("Setting election timeout at %v", randtval)
		electiontimer.Reset(randtval)
		<-electiontimer.C
		log.Printf("%d: calling election!\n", i)
		i++
	}

}
