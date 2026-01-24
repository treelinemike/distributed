// raftcommon.go
//
// Misc. shared types for Raft implementation.
//
// Author: M. Kokko
// Updated: 24-Jan-2025

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
