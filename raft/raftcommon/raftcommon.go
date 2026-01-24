package raftcommon

// state of raft server
type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)

type RespToClient struct {
	AppendedToLeader bool
	LeaderID         string
}
