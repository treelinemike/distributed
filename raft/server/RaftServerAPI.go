package main

import (
	"engg415/raft/common"
	"log"
)

type RaftAPI int

// TODO: struct members don't need to be exported if we're keeping in the server.go main package
type AEParams struct {
	Term         int
	LeaderId     string
	PrevLogIndex int
	LeaderCommit int
}

func (r *RaftAPI) AppendEntries(p AEParams, resp *int) error {
	log.Printf("AppendEntriesRPC received in term %d from leader %v\n", p.Term, p.LeaderId)
	electiontimer.Stop()
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
	log.Printf("Server %s requested a vote for term %v\n", p.CandidateId, p.Term)

	if currentTerm < p.Term {
		currentTerm = p.Term
		state = common.Follower
		votedForThisTerm = p.CandidateId
		resp.Term = currentTerm
		resp.VoteGranted = true
		log.Printf("Vote granted\n")
	} else if currentTerm == p.Term {
		resp.Term = currentTerm
		if votedForThisTerm == "" {
			votedForThisTerm = p.CandidateId
			resp.VoteGranted = true
			log.Printf("Vote granted\n")
		} else {
			resp.VoteGranted = false
			log.Printf("Vote denied\n")
		}
	}
	// todo: other important cases here?

	return nil
}
