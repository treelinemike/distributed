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
	if p.Term > currentTerm {
		log.Printf("Learned from server %v that we're actually in term %v (not %v)\n", p.LeaderID, p.Term, currentTerm)
		state = common.Follower
		currentTerm = p.Term
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

	if currentTerm < p.Term {
		currentTerm = p.Term
		state = common.Follower
		votedForThisTerm = p.CandidateId
		resp.Term = currentTerm
		resp.VoteGranted = true
		log.Printf("Term incremented. Vote granted for term %v\n", currentTerm)
	} else if currentTerm == p.Term {
		resp.Term = currentTerm
		if votedForThisTerm == "" { // not sure we should ever get to this case?
			votedForThisTerm = p.CandidateId
			resp.VoteGranted = true
			log.Printf("Vote granted for term %v\n", currentTerm)
		} else { // we've already voted for someone
			resp.VoteGranted = false
			log.Printf("Vote denied, we've already voted for %v in term %v\n", votedForThisTerm, currentTerm)
		}
	} else {
		resp.VoteGranted = false
		log.Printf("Vote denied. Candidate %v term (%v) is behind current term (%v).\n", p.CandidateId, p.Term, currentTerm)
	}

	// reset election timer here
	// raft paper doesn't appear to do this, but raft scope does
	// and it seems to help reduce number of unnecessary elections
	electiontimer.Reset()

	return nil
}

func (r *RaftAPI) ProcessClientRequest(s []string, resp *common.RespToClient) error {

	// in all cases, return current leader ID
	resp.LeaderID = currentTermLeader

	// if not leader, send client back the current leader ID
	if state != common.Leader {
		log.Printf("Redirecting a client request to leader (%v) to commit: %v\n", currentTermLeader, s)
		resp.Committed = false
		return nil
	}

	// if we are the leader, try to commit this to each server
	// TODO: DO THIS HERE
	log.Printf("Processing client request to commit: %v\n", s)
	resp.Committed = true
	return nil
}

func (r *RaftAPI) StopServer(param int, resp *int) error {
	log.Println("Received STOP from client, process will end on next loop")
	stopServer = true
	return nil
}
