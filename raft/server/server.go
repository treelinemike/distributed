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
	common.LoadRaftConfig(filename, servers, &t)
	fmt.Printf("Timeout range: [%d, %d]\n", t.Min_ms, t.Max_ms)
	fmt.Printf("Random draw: %d\n", rand.IntN(t.Max_ms-t.Min_ms)+t.Min_ms) // test picking a random timeout value

	// make sure provided selfkey is in map from config file
	_, ok := servers[selfkey]
	if !ok {
		log.Fatal("Key ", selfkey, " is not in cluster config file")
	}

	// serve RaftAPI
	log.Println("Registering server API for access on port", servers[selfkey].Port)
	thisapi := new(RaftAPI)
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":"+servers[selfkey].Port)
	if err != nil {
		log.Fatal("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)
	log.Println("Ready to play")

	// start a timer
	i := 0
	timer1 := time.NewTimer(2 * time.Second)
	for {
		<-timer1.C
		fmt.Printf("%d: calling election!\n", i)
		i++
		timer1.Reset(2 * time.Second)
	}

}
