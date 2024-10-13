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
	"sync"
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
	var jsonfilebase string
	common.LoadRaftConfig(filename, servers, &t, &jsonfilebase)
	log.Printf("Config specifies election timeout range [%d, %d]\n", t.Min_ms, t.Max_ms)

	// make sure provided selfkey is in map from config file
	_, ok := servers[selfkey]
	if !ok {
		log.Fatal("Key ", selfkey, " is not in cluster config file")
	} else {
		log.Printf("Config specifies this server as %s:%s", servers[selfkey].Address, servers[selfkey].Port)
	}

	// load non-volatile state if it has been previously saved
	jsonfilename = jsonfilebase + "_" + selfkey + ".json"
	log.Println("Attempting to load non-volatile state from file: ", jsonfilename)
	err = readnvstate()
	if err != nil {
		log.Fatal("Error initializing initial state")
	}

	// test setting term
	setterm(14)
	setleaderid("someserver")

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
	for svr := range servers {
		if svr != selfkey {
			wg.Add(1)
			go func(svr string) {
				defer wg.Done()
				// try connect to host repeatedly (even if OS times out the attempt very quickly)
				// we will use our own timeout on the connection attempt
				for {
					log.Printf("Attempting connection to host %s:%s\n", servers[svr].Address, servers[svr].Port)
					host, err := rpc.DialHTTP("tcp", servers[svr].Address+":"+servers[svr].Port)
					if err == nil {
						log.Printf("Connected to host %s:%s\n", servers[svr].Address, servers[svr].Port)

						// store handle to server in struct in map
						// Go doesn't make this easy
						tempsvrdata := servers[svr]
						tempsvrdata.Handle = host
						servers[svr] = tempsvrdata
						break
					}
				}
			}(svr)
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
	case <-time.After(20 * time.Second):
		log.Fatal("Timeout connecting to servers")
	}

	// start the election timeout timer
	electiontimer = time.NewTimer(100 * time.Second)
	i := 0
	for {
		randtval := time.Duration(rand.IntN(t.Max_ms-t.Min_ms)+t.Min_ms) * time.Millisecond
		log.Printf("Setting election timeout at %v", randtval)
		electiontimer.Reset(randtval)
		<-electiontimer.C
		log.Printf("%d: calling election!\n", i)
		i++

		// request votes from all other servers
		for svr := range servers {
			if svr != selfkey {
				log.Printf("We should ask server %s: %s:%s ", svr, servers[svr].Address, servers[svr].Port)
			}
		}

	}

}
