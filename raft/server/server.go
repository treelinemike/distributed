package main

import (
	"engg415/raft/raftcommon"
	"engg415/raft/raftnv"
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

// globals that need to be accessed by functions or RPCs

// non-volatile storage for current term, voted for, and log entries
var st raftnv.NVState

// volatile storage on this server
var currentTermLeader string = ""
var commitIdx int = 0
var nextIdx map[string]int = make(map[string]int)
var matchIdx map[string]int = make(map[string]int)
var state raftcommon.RaftState = raftcommon.Follower
var electiontimer raftcommon.RaftTimer
var servers map[string]raftcommon.NetworkAddress
var isactive map[string]bool
var stopServer bool = false

// function to reset nextIdx and matchIdx upon election as leader
func resetLeaderIndices() {
	for svr := range servers {
		if svr != "" {
			nextIdx[svr] = st.LogLength() + 1
			matchIdx[svr] = 0
		}
	}
}

// function to repeatedly attempt reconnecting to another server
// will be launched as a goroutine
func reconnectToHost(svr string) {
	log.Printf("Attempting to reconnect to server %v (%s:%s)\n", svr, servers[svr].Address, servers[svr].Port)
	for {
		host, err := rpc.DialHTTP("tcp", servers[svr].Address+":"+servers[svr].Port)
		if err == nil {
			log.Printf("Reconnected to server %v (%s:%s)\n", svr, servers[svr].Address, servers[svr].Port)

			// store handle to server in struct in map
			// Go doesn't make this easy
			tempsvrdata := servers[svr]
			tempsvrdata.Handle = host
			servers[svr] = tempsvrdata

			// update to indicate that server is back online
			isactive[svr] = true

			// break out of loop
			break
		}
	}
}

// main function/goroutine for the peer raft servers
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
	servers = make(map[string]raftcommon.NetworkAddress)
	isactive = make(map[string]bool)

	// load config
	log.Println("Loading cluster configuration file")
	var jsonfilebase string
	var t raftcommon.Timeout
	raftcommon.LoadRaftConfig(filename, servers, &t, &jsonfilebase)

	// create an election timer
	log.Printf("Config specifies election timeout range [%d, %d]\n", t.Min_ms, t.Max_ms)
	electiontimer.SetTimeout(t)
	electiontimer.Timer = time.NewTimer(100 * time.Second) // TODO: add this into an init function?

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
	st.SetJSONFilename(jsonfilebase + "_" + selfkey + ".json")
	log.Println("Attempting to load non-volatile state from file: ", st.GetJSONFilename())
	err = st.ReadNVState()
	if err != nil {
		log.Fatal("Error initializing initial state")
	}

	// show what we loaded from json
	log.Printf("Loaded voted for: '%v'\n", st.GetVotedFor())
	log.Printf("Loaded term: %v\n", st.GetCurrentTerm())
	log.Printf("Loaded log: %v\n", st.Log) // accessing field directly for debug/logging only

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
	// ok to block with a wait group here beause we do actually need to connect to all of our servers
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

						// set "up" status
						isactive[svr_key] = true

						// initialize nextIdx and matchIdx for this server
						nextIdx[svr_key] = st.LogLength() + 1
						matchIdx[svr_key] = 0

						// break out of loop
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
	state = raftcommon.Follower

	for !stopServer {

		// check to make sure all servers we expect to be active are still connected
		// TODO: abstract network layer to its own module/file
		for svr := range servers {
			if svr != selfkey && isactive[svr] {
				var retval int
				err := servers[svr].Handle.Call("RaftAPI.Ping", 1, &retval)
				if err != nil {
					log.Printf("Error pinging server %v: %v\n", svr, err)
					isactive[svr] = false
					go reconnectToHost(svr)
				}
			}
		}

		// handle server states (follower, candidate, leader)
		switch state {
		case raftcommon.Follower:
			electiontimer.Reset()
			<-electiontimer.Timer.C
			state = raftcommon.Candidate

		case raftcommon.Candidate:

			// increment term
			st.SetCurrentTerm(st.GetCurrentTerm() + 1)
			st.SetVotedFor(selfkey)
			st.WriteNVState()

			// vote for self
			numVotesReceived := 1
			numBallotsReceived := 1
			numBallotsSent := 1

			log.Printf("Calling election for term %v\n", st.GetCurrentTerm())
			electionTerm := st.GetCurrentTerm()

			// create a channel to receive votes
			votechan := make(chan bool)

			// request votes from all other servers
			// TODO: we're not retrying any servers that don't respond here (raftscope does show retries every ~50ms)
			for svr := range servers {
				if svr != selfkey && isactive[svr] {

					// launch a goroutine to request vote on a channel
					go func(svr string, votechannel chan bool) {
						log.Printf("Requesting term %v vote from server %s", st.GetCurrentTerm(), svr)

						// prepare parameters for requesting vote
						var rvp RVParams
						rvp.CandidateId = selfkey
						rvp.LastLogIndex = 0 // TODO: add correct value
						rvp.Term = st.GetCurrentTerm()

						var rvr RVResp
						err := servers[svr].Handle.Call("RaftAPI.RequestVote", rvp, &rvr)
						// if server has been shutdown, try recommencting once per RPC attempt
						// TODO: move this to an ASYNCHRONOUS function, will require making 'servers' a global
						if err == rpc.ErrShutdown {
							isactive[svr] = false
							go reconnectToHost(svr)
						}
						if err != nil {
							log.Printf("Could not request term %v vote from server %s: %v\n", st.CurrentTerm, svr, err)
						}

						votechannel <- rvr.VoteGranted // defaults to false
					}(svr, votechan)
					numBallotsSent += 1

				}
			}

			// run election timer
			electiontimer.Reset()

			for done := false; !done; {
				select {
				case vote := <-votechan:

					// log that we got a ballot
					numBallotsReceived += 1

					// determne whether we got the vote
					if vote {
						numVotesReceived += 1
						log.Printf("Ballot received: vote granted")
					} else {
						log.Printf("Ballot received: vote denied")
					}

					// exit if we received all of the ballots
					// TODO: this is a bit of a hack to avoid goroutines writing to a closed channel
					// if we don't get all the ballots back we will therefore wait until the election timer expires
					// to declare victory even if we receive a majority of the votes much earlier
					if numBallotsReceived == numBallotsSent {
						done = true
						if st.GetCurrentTerm() == electionTerm && float32(numVotesReceived) > float32(len(servers))/2.0 { // need a strict inequality here (true majority)
							log.Printf("We won the term %v election\n", st.GetCurrentTerm())
							state = raftcommon.Leader
							resetLeaderIndices()
						} else if st.GetCurrentTerm() != electionTerm {
							log.Printf("Term changed during election to %v, reverting to follower\n", st.GetCurrentTerm())
							state = raftcommon.Follower
							st.SetVotedFor("")
						} else {
							log.Printf("We did not win the term %v election\n", st.GetCurrentTerm())
							state = raftcommon.Follower
							st.SetVotedFor("")
						}
					}

				// timeout case
				case <-electiontimer.Timer.C:
					log.Printf("Term %v election timed out waiting for votes\n", st.GetCurrentTerm())
					done = true

					// check to see if we won...
					// TODO: remove this once we fix the logic in the election piece (i.e. win immediately on majority)
					if st.GetCurrentTerm() == electionTerm && float32(numVotesReceived) > float32(len(servers))/2.0 { // need a strict inequality here (true majority)
						log.Printf("We won the term %v election but only received %v ballots\n", st.GetCurrentTerm(), numBallotsReceived)
						state = raftcommon.Leader
						resetLeaderIndices()
					} else if st.GetCurrentTerm() != electionTerm {
						log.Printf("Term changed during election to %v, reverting to follower\n", st.GetCurrentTerm())
						state = raftcommon.Follower
						st.SetVotedFor("")
					}

					// otherwise we remain a candidate and will re-run the election

				}
			}

			// report election statistics for the log
			log.Printf("Received %v votes in %v ballots\n", numVotesReceived, numBallotsReceived)

			// don't change or re-initiate our canddidate state
			// because it may be changed by incoming RPCs (e.g. if we have the wrong term)

			// close votechan
			close(votechan)

		case raftcommon.Leader:

			// sent AppendEntries RPCs to all followers
			for svr := range servers {
				if svr != selfkey && !isactive[svr] {
					log.Printf("NOT sending term %v AppendEntries to offline server %v\n", st.CurrentTerm, svr)
				}
				if svr != selfkey && isactive[svr] {
					log.Printf("Sending term %v AppendEntries to server %s\n", st.CurrentTerm, svr)

					// prepare parameters for calling append entries RPC
					var aeresp AEResp
					var aeparam AEParams
					aeparam.Term = st.GetCurrentTerm()
					aeparam.LeaderID = selfkey
					aeparam.LeaderCommit = commitIdx
					aeparam.PrevLogIndex = nextIdx[svr] - 1   // previous log intex always one before next index
					aeparam.Entries = []raftnv.RaftLogEntry{} // default: empty slice of entries, this would happen anyway

					// prepare parameters for AppendEntries RPC based on follower status
					if nextIdx[svr] == (st.LogLength() + 1) {

						// follower log is up to date -> send heartbeat (empty) AERPC
						// this case should also run when both leader and follower logs are empty
						aeparam.PrevLogTerm = 0
						if st.LogLength() > 0 {
							prevEntry, err := st.GetLogEntry(st.LogLength())
							if err != nil {
								log.Panicf("Error getting previous log entry for server %s: %v\n", svr, err)
							}
							aeparam.PrevLogTerm = prevEntry.Term
						}

					} else if nextIdx[svr] == 1 {

						// follower log is empty -> sent first entry
						// the heartbeat case coming first means we know there is at least one entry in the leader log
						aeparam.PrevLogTerm = 0
						nextEntry, err := st.GetLogEntry(nextIdx[svr])
						if err != nil {
							log.Panicf("Error getting next log entry for server %s: %v\n", svr, err)
						}
						aeparam.Entries = []raftnv.RaftLogEntry{nextEntry}

					} else if nextIdx[svr] <= st.LogLength() {

						// follower log is behind leader -> send next entry
						prevEntry, err := st.GetLogEntry(nextIdx[svr] - 1)
						if err != nil {
							log.Panicf("Error getting previous log entry for server %s: %v\n", svr, err)
						}
						nextEntry, err := st.GetLogEntry(nextIdx[svr])
						if err != nil {
							log.Panicf("Error getting next log entry for server %s: %v\n", svr, err)
						}
						aeparam.PrevLogTerm = prevEntry.Term
						aeparam.Entries = []raftnv.RaftLogEntry{nextEntry}

					} else {

						// that SHOULD be all of our cases...
						log.Panicf("Invalid nextIdx case planning AppendEntries call\n")
					}

					// SEND IT!
					err := servers[svr].Handle.Call("RaftAPI.AppendEntries", aeparam, &aeresp)
					if err != nil {
						log.Printf("Error calling AppendEntries on server %s: %v\n", svr, err)
					}

					// if server has been shutdown, try recommencting once per RPC attempt
					// TODO: move this to an ASYNCHRONOUS function, will require making 'servers' a global
					if err == rpc.ErrShutdown {
						isactive[svr] = false
						go reconnectToHost(svr)
						continue
					}

					// if the follower is at a more current term we need to step down
					if !aeresp.Success && aeresp.Term > st.GetCurrentTerm() {
						log.Printf("Follower %s term %v greater than our term %v\n", svr, aeresp.Term, st.GetCurrentTerm())
						log.Printf("Stepping down to follower\n")
						state = raftcommon.Follower
						st.SetCurrentTerm(aeresp.Term)
						st.SetVotedFor("")
						st.WriteNVState()
						continue
					}

					// if AppendtEntries RPC failed, decrement nextIdx for that server (will try again next time)
					if !aeresp.Success {
						log.Printf("AppendEntries RPC to server %s FAILED\n", svr)
						log.Printf("Decrement nextIdx %v -> %v\n", nextIdx[svr], nextIdx[svr]-1)
						nextIdx[svr] -= 1
					} else {
						// success case, increment nextIdx and matchIdx for that server
						log.Printf("AppendEntries RPC to server %s SUCCEEDED\n", svr)
						if len(aeparam.Entries) > 0 {
							log.Printf("Increment nextIdx %v -> %v", nextIdx[svr], nextIdx[svr]+len(aeparam.Entries))
							nextIdx[svr] += len(aeparam.Entries)
							log.Printf("Update matchIdx %v -> %v\n", matchIdx[svr], nextIdx[svr]+len(aeparam.Entries)-1)
							matchIdx[svr] = nextIdx[svr] - 1 // TODO: check this logic, not clear in Raft paper?
						}
					}

				}
			}

			// determine whether any new entries can be committed
			// TODO: paper also has condition that entry must be from current term...
			for i := commitIdx + 1; i <= st.LogLength(); i++ {
				// count how many servers have this log entry
				count := 0
				for _, match := range matchIdx {
					if match >= i {
						count++
					}
				}
				// if majority of servers have this log entry, commit it
				// note: we include ourselves (the leader) in the count here
				if (count + 1) > (len(servers)+1)/2 {
					entr, err := st.GetLogEntry(i)
					if err != nil {
						log.Panicf("Error getting entry %v: %v\n", i, err)
					}
					log.Printf("Committing entry at index %v [%v])\n", i, entr)
					commitIdx = i
				}
			}

			// wait between AppendEntries RPCs
			// TODO: if we have something to send we should send it immediately rather than waiting for the heartbeat interval
			time.Sleep(500 * time.Millisecond)

		default:
			log.Fatal("Unknown state: ", state)

		}

	}

}
