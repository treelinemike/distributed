package common

// state of raft server
type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)
