package main

import (
	"engg415/raft/common"
	"log"
)

type RaftAPI int

// TODO: struct members don't need to be exported if we're keeping in the server.go main package
type AEParams struct {
	Term         int
	LeaderID     string
	PrevLogIndex int
	LeaderCommit int
}

func (r *RaftAPI) AppendEntries(p AEParams, resp *int) error {
	log.Printf("AppendEntriesRPC received in term %d from leader %v\n", p.Term, p.LeaderID)

	// update LeaderID so we can redirect any client requests
	currentTermLeader = p.LeaderID

	// stop trying to win an election (and importantly incrementing the term to do so)
	if p.Term > st.CurrentTerm {
		log.Printf("Learned from server %v that we're actually in term %v (not %v)\n", p.LeaderID, p.Term, st.CurrentTerm)
		state = common.Follower
		st.CurrentTerm = p.Term
		st.VotedFor = ""
		writenvstate()
	}
	electiontimer.Reset()

	// TODO: handle log consistency and synchronization here...

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
	log.Printf("Server %s requested a vote for term %v\n", p.CandidateId, p.Term)

	if st.CurrentTerm < p.Term {
		st.CurrentTerm = p.Term
		state = common.Follower
		st.VotedFor = p.CandidateId
		writenvstate()
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
		log.Printf("Vote denied. Candidate %v term (%v) is behind current term (%v).\n", p.CandidateId, p.Term, st.CurrentTerm)
	}

	// reset election timer here
	// raft paper doesn't appear to do this, but raft scope does
	// and it seems to help reduce number of unnecessary elections
	electiontimer.Reset()

	return nil
}

// RPC for clients to submit requests to be committed
func (r *RaftAPI) ProcessClientRequest(s []string, resp *common.RespToClient) error {

	// in all cases, return current leader ID
	resp.LeaderID = currentTermLeader
	resp.AppendedToLeader = false

	// if not leader, send client back the current leader ID
	if state != common.Leader {
		log.Printf("Redirecting a client request to leader (%v) to commit: %v\n", currentTermLeader, s)
		return nil
	}

	// if we are the leader, add to our log which will trigger replication
	log.Printf("Processing client request to commit: %v\n", s)
	for _, v := range s {
		common.AppendLogEntry(st.CurrentTerm, v)
	}
	resp.AppendedToLeader = true
	return nil
}

// RPC to allow clients to check whether all requests are committed
func (r *RaftAPI) IsFullyCommitted(param int, resp *bool) error {
	*resp = (commitIdx == len(common.RaftLog)) // remember commitIdx is 1-based
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
