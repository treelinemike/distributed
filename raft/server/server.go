package main

import (
	"engg415/raft/common"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type RaftAPI int

// globals that need to be accessed by RPCs
var currentTerm int = 0
var votedForThisTerm string = ""
var commitIdx int = 0
var lastApplied int = 0
var state common.RaftState = common.Follower
var electiontimer common.RaftTimer

// TODO: struct members don't need to be exported if we're keeping in the server.go main package
type AEParams struct {
	term         int
	leaderId     int
	prevLogIndex int
	leaderCommit int
}

func (r *RaftAPI) AppendEntries(p AEParams, resp *int) error {
	electiontimer.Stop()
	fmt.Printf("AppendEntriesRPC: term %d\n", p.term)

	electiontimer.Reset()
	return nil
}

type RVParams struct {
	Term         int
	CandidateId  string
	LastLogIndex int
	LastLogTerm  int
}

type RVResp struct {
	Term        int
	VoteGranted bool
}

func (r *RaftAPI) RequestVote(p RVParams, resp *RVResp) error {
	fmt.Printf("RequestVoteRPC: term %d vote requested for server %s\n", p.Term, p.CandidateId)

	return nil
}

func main() {

	var err error

	// format log time in microseconds
	// will start logging to file as soon as we know which host key we have been assigned
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
	log.Println("Loading cluster configuration file")
	var jsonfilebase string
	var t common.Timeout
	common.LoadRaftConfig(filename, servers, &t, &jsonfilebase)
	log.Printf("Config specifies election timeout range [%d, %d]\n", t.Min_ms, t.Max_ms)
	electiontimer.SetTimeout(t)

	// make sure provided selfkey is in map from config file
	_, ok := servers[selfkey]
	if !ok {
		log.Fatal("Key ", selfkey, " is not in cluster config file")
	} else {
		log.Printf("Config specifies this server as %s:%s with host key %s", servers[selfkey].Address, servers[selfkey].Port, selfkey)
	}

	// configure logging
	logfilename := fmt.Sprintf("log_%s.txt", selfkey)
	log.Printf("Creating log file: %s", logfilename)
	logfid, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Panicf("Could not create log file: %v\n", err)
	}
	defer logfid.Close()
	mux := io.MultiWriter(os.Stdout, logfid)
	log.SetOutput(mux)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) // for good measure?

	// load non-volatile state if it has been previously saved
	jsonfilename = jsonfilebase + "_" + selfkey + ".json"
	log.Println("Attempting to load non-volatile state from file: ", jsonfilename)
	err = readnvstate()
	if err != nil {
		log.Fatal("Error initializing initial state")
	}

	// show what we loaded from json
	log.Printf("Loaded leaderID: '%v'\n", st.LeaderID)
	log.Printf("Loaded term: %v\n", st.Term)
	log.Printf("Loaded log: %v\n", st.Log)

	// serve up our RPC API
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

	// connect to each server in the cluster
	// use a bunch of goroutines with timeouts
	// TODO: better handle the case when a server is initially crashed?
	log.Println("Trying to connect to all other servers now...")
	var wg sync.WaitGroup
	for svr_key, svr_data := range servers { // ranging over a map returns (key, value) pairs
		if svr_key != selfkey {
			wg.Add(1)
			go func(svr_key string) { // note: we can still access servers due to scoping of goroutine
				defer wg.Done()
				// try connect to host repeatedly (even if OS times out the attempt very quickly)
				// we will use our own timeout on the connection attempt
				for {
					//log.Printf("Attempting connection to host %s:%s\n", servers[svr].Address, servers[svr].Port)
					host, err := rpc.DialHTTP("tcp", servers[svr_key].Address+":"+servers[svr_key].Port)
					if err == nil {
						log.Printf("Connected to host %s:%s\n", servers[svr_key].Address, servers[svr_key].Port)

						// store handle to server in struct in map
						// Go doesn't make this easy
						tempsvrdata := svr_data
						tempsvrdata.Handle = host
						servers[svr_key] = tempsvrdata
						break
					}
				}
			}(svr_key)
		}
	}
	c := make(chan struct{})
	go func() {
		wg.Wait()
		defer close(c)
	}()
	select {
	case <-c:
		log.Println("Connected to all servers")
	case <-time.After(5 * time.Minute):
		log.Fatal("Timeout connecting to servers")
	}

	// by default start in follower state
	log.Println("Starting in follower state")
	state = common.Follower

	if selfkey != "1" {
		for {

		}
	}

	// start the election timeout timer
	electiontimer.Timer = time.NewTimer(100 * time.Second)

	i := 0
	for {
		log.Printf("Starting election timer\n")
		electiontimer.Reset()
		<-electiontimer.Timer.C
		log.Printf("%d: calling election!\n", i)
		i++

		// prepare parameters for requesting vote
		var rvp RVParams
		rvp.CandidateId = selfkey
		rvp.LastLogIndex = 0 // TODO: add correct value
		rvp.Term = 0         // TODO: add correct value

		// request votes from all other servers
		for svr := range servers {
			if svr != selfkey {
				log.Printf("We should ask server %s: %s:%s ", svr, servers[svr].Address, servers[svr].Port)
				var rvr RVResp
				err := servers[svr].Handle.Call("RaftAPI.RequestVote", rvp, &rvr)
				if err != nil {
					log.Fatalf("Could not request vote from server %s: %v\n", svr, err)
				}
			}
		}

	}

}
