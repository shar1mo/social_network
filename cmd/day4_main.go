package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

// Роль узла
type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

// Тип сообщения
type Message struct {
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
	resetTimer  chan struct{}
	Alive       bool
	countVotes  int
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
	Inbox       chan Message
	mutex       sync.Mutex
}

func InitializeNodes(ids []int) map[int]*Node {
	// Инициализируем узлы
	var my_role Role = Follower
	result := map[int]*Node{}
	for _, id := range ids {
		result[id] = &Node{Alive: true, ID: id, State: my_role, LeaderID: -1, Peers: filter(ids, id), VotedFor: -1, Inbox: make(chan Message, 100)}
	}

	return result
}

// Фильтрует список узлов, исключая текущий
func filter(ids []int, exclude int) []int {
	peers := []int{}
	for _, id := range ids {
		if id != exclude {
			peers = append(peers, id)
		}
	}
	return peers
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
		case "LeaderApply":
			n.handleLeaderApply(msg, nodes)
		case "HeartBeat":
			n.handleHeartBeats(msg, nodes)
		}
	}
}

func (n *Node) runElectionTimer(nodes map[int]*Node) {
	timeout := time.Duration(rand.Intn(200)+200) * time.Millisecond
	electionTimer := time.NewTimer(timeout)
	for {
		select {
		case <-electionTimer.C:
			fmt.Printf("%d: My timer stopped. Start election\n", n.ID)
			if n.State == 0 {
				print("Aboba")
				n.State = Candidate // Переход в состояние кандидата
				n.CurrentTerm++     // Увеличиваем текущий термин
				n.VotedFor = n.ID   // Проголосовать за себя
				// Запуск выборов
				n.Inbox <- Message{Type: "RequestVoteReply", FromID: n.ID}
			}

			timeout = time.Duration(rand.Intn(150)+150) * time.Millisecond
			electionTimer.Reset(timeout)

		case <-n.resetTimer:
			// Сброс таймера при получении сигнала
			fmt.Printf("%d: Resetting election timer\n", n.ID)
			electionTimer.Reset(timeout)
			if n.State != Leader { // Сброс таймера только если не лидер
				fmt.Printf("%d: Resetting election timer\n", n.ID)
				electionTimer.Reset(timeout)
			} else {
				fmt.Printf("%d: I am the leader, not resetting timer\n", n.ID)
			}
			n.sendHeartbeats(nodes)
		}
	}
}

func (n *Node) sendHeartbeats(nodes map[int]*Node) {
	if n.State == Leader {
		for _, peer := range n.Peers {
			fmt.Printf("%d: Send HeartBeat to %d\n", n.ID, peer)
			nodes[peer].Inbox <- Message{
				Type:   "HeartBeat",
				FromID: n.ID,
				ToID:   peer,
				Term:   n.CurrentTerm,
			}
		}
	}
}

func (n *Node) handleHeartBeats(msg Message, nodes map[int]*Node) {
	if n.Alive {
		fmt.Printf("%d: Get HeartBeat from %d. Send OK\n", n.ID, msg.FromID)
		n.resetTimer <- struct{}{} // Сброс таймера
		nodes[msg.FromID].Inbox <- Message{
			Type:   "OK",
			FromID: n.ID,
			ToID:   msg.FromID,
		}
	}
}

func (n *Node) handleLeaderApply(msg Message, nodes map[int]*Node) {
	if n.Alive {
		fmt.Printf("%d: Get LeaderApply from %d. Update my LeaderID\n", n.ID, msg.FromID)
		n.LeaderID = msg.FromID
		n.resetTimer <- struct{}{}
	}
}

func (n *Node) handleRequestVote(msg Message, nodes map[int]*Node) {
	if n.CurrentTerm < msg.Term && n.VotedFor == -1 {
		fmt.Printf("%d: Get RequestVote from %d. Update my CurrentTemp and cast my vote\n", n.ID, msg.FromID)
		n.CurrentTerm = msg.Term
		n.VotedFor = msg.FromID
		nodes[msg.FromID].Inbox <- Message{Type: "RequestVoteReply", FromID: n.ID, ToID: msg.FromID, VoteGranted: true}
	} else {
		fmt.Printf("%d: Get RequestVote from %d. Dont cast my vote\n", n.ID, msg.FromID)

	}

}

func (n *Node) handleRequestVoteReply(msg Message, nodes map[int]*Node) {
	if n.countVotes == 0 {
		fmt.Printf("%d: Get RequestVotesReply from %d\n", n.ID, msg.FromID)
		n.countVotes++
		for _, peer := range n.Peers {
			fmt.Printf("%d: Send RequestVote to %d\n", n.ID, nodes[peer].ID)
			nodes[peer].Inbox <- Message{Type: "RequestVote", Term: n.CurrentTerm, FromID: n.ID, ToID: nodes[peer].ID}
		}
	} else {
		if msg.VoteGranted {
			fmt.Printf("%d: Get Vote from %d. Update my CountVotes\n", n.ID, msg.FromID)
			n.countVotes++
		}

		if n.countVotes >= (len(nodes)-1)/2 {
			fmt.Printf("%d: My CountVotes more then half of all nodes. Reset my Votes\n", n.ID, msg.FromID)
			n.State = Leader
			n.countVotes = 0
			for _, peer := range n.Peers {
				fmt.Printf("%d: Send LeaderApply to %d\n", n.ID, nodes[peer].ID)
				nodes[peer].Inbox <- Message{Type: "LeaderApply", FromID: n.ID, ToID: nodes[peer].ID}
			}

			go n.heartbeatLoop(nodes)

		}
	}
}

func (n *Node) heartbeatLoop(nodes map[int]*Node) {
	for n.State == Leader {
		n.sendHeartbeats(nodes)
		//time.Sleep(100 * time.Millisecond) // Периодичность отправки Heartbeat
	}
}

func main() {
	ids := []int{1, 2, 3}
	nodes := InitializeNodes(ids)
	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go node.Run(&wg, nodes)
	}

	time.Sleep(4 * time.Second)

	for _, node := range nodes {
		close(node.Inbox)
	}

	wg.Wait()

}
