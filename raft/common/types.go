package common

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

type RaftLogEntry struct {
	Term  int
	Value string
}

var RaftLog []RaftLogEntry = make([]RaftLogEntry, 0)

func AppendLogEntry(term int, value string) {
	entry := RaftLogEntry{
		Term:  term,
		Value: value,
	}
	RaftLog = append(RaftLog, entry)
}

func GetLogTerm(index int) int {
	if index < 1 || index > len(RaftLog) {
		return -1
	}
	return RaftLog[index-1].Term
}
