package distributed

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Роль узла
type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

// Тип сообщения
type MessageRaft struct {
	Term         int
	FromID       int
	ToID         int
	Type         string // "RequestVote", "RequestVoteReply", "AppendEntries", "AppendEntriesReply"
	LastLogIndex int
	LastLogTerm  int
	VoteGranted  bool
	Success      bool
	Entries      []LogEntry
	LeaderCommit int
}

// Запись лога
type LogEntry struct {
	Term    int
	Command interface{}
}

// Структура узла
type Node struct {
	ID          int
	LeaderID    int
	Peers       []int
	State       Role
	CurrentTerm int
	VotedFor    int
	Log         []LogEntry
	CommitIndex int
	LastApplied int
	NextIndex   map[int]int
	MatchIndex  map[int]int
	Inbox       chan MessageRaft
	countVotes  int
	mutex       sync.Mutex
}

// Фильтрует список узлов, исключая текущий
func filter(ids []int, exclude int) []int {
	var peers []int
	for _, id := range ids {
		if id != exclude {
			peers = append(peers, id)
		}
	}
	return peers
}

// Инициализация узлов
func InitializeNodes(ids []int) map[int]*Node {
	result := make(map[int]*Node)
	for _, id := range ids {
		result[id] = &Node{
			ID:         id,
			State:      Follower,
			LeaderID:   -1,
			Peers:      filter(ids, id),
			VotedFor:   -1,
			Inbox:      make(chan MessageRaft, 10),
			NextIndex:  make(map[int]int),
			MatchIndex: make(map[int]int),
		}
	}
	return result
}

func (n *Node) Run(wg *sync.WaitGroup, nodes map[int]*Node) {
	defer wg.Done()
	go n.runElectionTimer(nodes)

	for msg := range n.Inbox {
		switch msg.Type {
		case "RequestVote":
			n.handleRequestVote(msg, nodes)
		case "RequestVoteReply":
			n.handleRequestVoteReply(msg, nodes)
		case "AppendEntries":
			n.handleAppendEntries(msg, nodes)
		case "AppendEntriesReply":
			n.handleAppendEntriesReply(msg, nodes)
		}
	}
}

func (n *Node) runElectionTimer(nodes map[int]*Node) {
	for {
		timeout := time.Duration(rand.Intn(150)+150) * time.Millisecond
		electionTimer := time.NewTimer(timeout)

		<-electionTimer.C

		n.mutex.Lock()
		if n.State == Follower {
			fmt.Printf("Node %d: Election timeout, becoming Candidate\n", n.ID)
			n.State = Candidate
			n.CurrentTerm++
			n.VotedFor = n.ID
			n.countVotes = 1 // Голосуем за себя

			for _, peerID := range n.Peers {
				nodes[peerID].Inbox <- MessageRaft{
					Type:         "RequestVote",
					Term:         n.CurrentTerm,
					FromID:       n.ID,
					LastLogIndex: len(n.Log) - 1,
					LastLogTerm:  n.getLastLogTerm(),
				}
			}
		}
		n.mutex.Unlock()
	}
}

func (n *Node) handleRequestVote(msg MessageRaft, nodes map[int]*Node) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if msg.Term < n.CurrentTerm {
		return
	}

	if msg.Term > n.CurrentTerm {
		n.CurrentTerm = msg.Term
		n.State = Follower
		n.VotedFor = -1
	}

	if n.VotedFor == -1 || n.VotedFor == msg.FromID {
		n.VotedFor = msg.FromID
		fmt.Printf("Node %d: Voting for Node %d\n", n.ID, msg.FromID)
		nodes[msg.FromID].Inbox <- MessageRaft{
			Type:        "RequestVoteReply",
			FromID:      n.ID,
			ToID:        msg.FromID,
			Term:        n.CurrentTerm,
			VoteGranted: true,
		}
	}
}

func (n *Node) handleRequestVoteReply(msg MessageRaft, nodes map[int]*Node) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.State != Candidate || msg.Term < n.CurrentTerm {
		return
	}

	if msg.VoteGranted {
		n.countVotes++
		fmt.Printf("Node %d: Received vote from %d, total votes: %d\n", n.ID, msg.FromID, n.countVotes)
	}

	if n.countVotes > len(nodes)/2 {
		fmt.Printf("Node %d: Became Leader\n", n.ID)
		n.State = Leader
		n.LeaderID = n.ID

		for _, peerID := range n.Peers {
			n.NextIndex[peerID] = len(n.Log)
			n.MatchIndex[peerID] = 0
		}

		n.sendHeartbeats(nodes)
	}
}

func (n *Node) sendHeartbeats(nodes map[int]*Node) {
	for _, peerID := range n.Peers {
		nodes[peerID].Inbox <- MessageRaft{
			Type:   "AppendEntries",
			Term:   n.CurrentTerm,
			FromID: n.ID,
		}
	}
}

func (n *Node) handleAppendEntries(msg MessageRaft, nodes map[int]*Node) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if msg.Term < n.CurrentTerm {
		nodes[msg.FromID].Inbox <- MessageRaft{
			Type:    "AppendEntriesReply",
			Term:    n.CurrentTerm,
			Success: false,
		}
		return
	}

	n.CurrentTerm = msg.Term
	n.State = Follower
	n.LeaderID = msg.FromID

	nodes[msg.FromID].Inbox <- MessageRaft{
		Type:    "AppendEntriesReply",
		Term:    n.CurrentTerm,
		Success: true,
	}
}

func (n *Node) handleAppendEntriesReply(msg MessageRaft, nodes map[int]*Node) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if msg.Term > n.CurrentTerm {
		n.CurrentTerm = msg.Term
		n.State = Follower
		n.VotedFor = -1
		return
	}

	if msg.Success {
		n.MatchIndex[msg.FromID] = msg.LastLogIndex
		n.NextIndex[msg.FromID] = msg.LastLogIndex + 1
	} else {
		if n.NextIndex[msg.FromID] > 0 {
			n.NextIndex[msg.FromID]--
		}
	}
}

func (n *Node) getLastLogTerm() int {
	if len(n.Log) == 0 {
		return 0
	}
	return n.Log[len(n.Log)-1].Term
}
