package main

import (
	"engg415/raft/raftcommon"
	"engg415/raft/raftnv"
	"log"
)

type RaftAPI int

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

// TODO: struct members don't need to be exported if we're keeping in the server.go main package
type AEParams struct {
	Term         int
	LeaderID     string
	LeaderCommit int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []raftnv.RaftLogEntry // will be packaged into RaftLogEntry structs when appended to log
}

type AEResp struct {
	Term    int
	Success bool
}

func (r *RaftAPI) AppendEntries(p AEParams, resp *AEResp) error {
	log.Printf("AppendEntries received in term %d from leader %v\n", p.Term, p.LeaderID)
	defer electiontimer.Reset() // reset election timer whenever AppendEntries RPC is received

	// by default report failure
	resp.Success = false

	// update LeaderID so we can redirect any client requests
	currentTermLeader = p.LeaderID

	// reply false if leader term is behind our current term (per Fig 2)
	// leader needs to step down
	if p.Term < st.GetCurrentTerm() {
		log.Printf("Leader %v term %v is behind our current term %v\n", p.LeaderID, p.Term, st.GetCurrentTerm())
		log.Printf("Responding FALSE to AppendEntries\n")
		return nil
	}

	// update term as needed
	// if we're behind we should stop trying to win an election (and importantly incrementing the term to do so)
	// but in all cases proceed to check log consistency
	if p.Term > st.GetCurrentTerm() {
		log.Printf("Term update from leader %v (%v -> %v)\n", p.LeaderID, st.GetCurrentTerm(), p.Term)
		state = raftcommon.Follower
		st.SetCurrentTerm(p.Term)
		st.SetVotedFor("")
		st.WriteNVState()
	}

	// now that we're in the correct term, add to our response
	resp.Term = st.GetCurrentTerm()

	// if LeaderCommit > commitIdx, set commitIdx = min(LeaderCommit, index of last new entry)
	// needs to happen before heartbeat case is handled because heartbeat will just return
	if p.LeaderCommit > commitIdx {
		newCommitIdx := min(p.LeaderCommit, st.LogLength())
		log.Printf("LeaderCommit %v is greater than our commitIdx %v\n", p.LeaderCommit, commitIdx)
		log.Printf("Updating our commitIdx to %v\n", newCommitIdx)
		commitIdx = newCommitIdx
	}

	// deal with heartbeat case
	if len(p.Entries) == 0 {
		resp.Success = true
		log.Printf("Responding TRUE to AppendEntries (heartbeat)\n")
		return nil
	}

	// reply false if log isn't long enough to contain an entry at PrevLogIndex
	if p.PrevLogIndex > st.LogLength() {
		log.Printf("Log length (%v) is shorter than PrevLogIndex (%v)\n", st.LogLength(), p.PrevLogIndex)
		log.Printf("Responding FALSE to AppendEntries\n")
		return nil
	}

	// reply false if log doesn't contain an entry at PrevLogIndex whose term matches PrevLogTerm
	// skip this check if PrevLogIndex is 0 (meaning empty log)
	if p.PrevLogIndex > 0 {
		prevLogEntry, err := st.GetLogEntry(p.PrevLogIndex)
		if err != nil {
			log.Printf("Error getting log entry %v\n", p.PrevLogIndex)
			log.Printf("Responding FALSE to AppendEntries\n")
			return nil
		}
		if prevLogEntry.Term != p.PrevLogTerm {
			log.Printf("PrevLogIndex %v has term %v, but leader PrevLogTerm is %v\n", p.PrevLogIndex, prevLogEntry.Term, p.PrevLogTerm)
			log.Printf("Responding FALSE to AppendEntries\n")
			return nil
		}
	}

	// if an existing entry conflicts with a new one (same index but different terms), delete the existing entry and all that follow it
	// TODO: we're just going to append all entries from PrevLogIndex+1 onwards for simplicity
	st.PruneLog(p.PrevLogIndex + 1)

	// append any new entries not already in the log
	// TODO: any errors to handle here?
	for _, entry := range p.Entries {
		st.AppendLogEntry(entry.Term, entry.Value)
		log.Printf("Appended new log entry from leader [%v]\n", entry)
	}

	// if we haven't rejected the RPC yet then we should exit here reporting success
	log.Printf("Responding TRUE to AppendEntries\n")
	resp.Success = true
	return nil
}

func (r *RaftAPI) RequestVote(p RVParams, resp *RVResp) error {

	log.Printf("Server %s requested a vote for term %v\n", p.CandidateId, p.Term)

	if st.CurrentTerm < p.Term {
		st.CurrentTerm = p.Term
		state = raftcommon.Follower
		st.VotedFor = p.CandidateId
		st.WriteNVState()
		resp.Term = st.CurrentTerm
		resp.VoteGranted = true
		log.Printf("Term incremented. Vote granted for term %v\n", st.CurrentTerm)
	} else if st.CurrentTerm == p.Term {
		resp.Term = st.CurrentTerm
		if st.VotedFor == "" { // not sure we should ever get to this case?
			st.VotedFor = p.CandidateId
			resp.VoteGranted = true
			log.Printf("Vote granted for term %v\n", st.CurrentTerm)
		} else { // we've already voted for someone
			resp.VoteGranted = false
			log.Printf("Vote denied, we've already voted for %v in term %v\n", st.VotedFor, st.CurrentTerm)
		}
	} else {
		resp.VoteGranted = false
		log.Printf("Vote denied. Candidate %v term (%v) is behind our current term (%v)\n", p.CandidateId, p.Term, st.CurrentTerm)
	}

	// reset election timer here
	// raft paper doesn't appear to do this, but raft scope does
	// and it seems to help reduce number of unnecessary elections
	electiontimer.Reset()

	return nil
}

// RPC for clients to submit requests to be committed
func (r *RaftAPI) ProcessClientRequest(s []string, resp *raftcommon.RespToClient) error {

	// in all cases, return current leader ID
	resp.LeaderID = currentTermLeader
	resp.AppendedToLeader = false

	// if not leader, send client back the current leader ID
	if state != raftcommon.Leader {
		log.Printf("Client request to commit [%v] -> redirect to leader %v\n", s, currentTermLeader)
		return nil
	}

	// if we are the leader, add to our log which will trigger replication
	log.Printf("Processing client request to commit [%v]\n", s)
	for _, v := range s {
		st.AppendLogEntry(st.CurrentTerm, v)
	}
	resp.AppendedToLeader = true
	return nil
}

// RPC to allow clients to check whether all requests are committed
func (r *RaftAPI) IsFullyCommitted(param int, resp *bool) error {
	*resp = (commitIdx == len(st.Log)) // remember commitIdx is 1-based
	return nil
}

// RPC for clients to stop the server
func (r *RaftAPI) StopServer(param int, resp *int) error {
	log.Println("Received STOP from client, process will end on next loop")
	stopServer = true
	return nil
}

// simple ping function to test connectivity
// could be integrated into AppendEntries but this is simpler for testing
// could also just use IsFullyCommitted() but this is clearer
func (r *RaftAPI) Ping(param int, resp *int) error {
	*resp = 1
	return nil
}
