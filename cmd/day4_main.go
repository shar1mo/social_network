package main

import "sync"

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
	countVotes	int
	ID          int
	LeaderID	int
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
	result := map[int]*Node
	for _, id := ids{
		result[node] = Node{ID: id, State: Role{Follower}, LeaderID: -1, Peers: filter(ids, id), VotedFor: -1}
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
        case "AppendEntries":
            n.handleAppendEntries(msg, nodes)
        case "AppendEntriesReply":
            n.handleAppendEntriesReply(msg, nodes)
        }
    }
}

func (n *Node) runElectionTimer(nodes map[int]*Node) {
    timeout := time.Duration(rand.Intn(150)+150) * time.Millisecond 
    electionTimer := time.NewTimer(timeout)
	for {
        select {
        case <-electionTimer.C:
            n.mu.Lock()
            if n.State == Follower {
                n.state = Candidate // Переход в состояние кандидата
                n.CurrentTerm++       // Увеличиваем текущий термин
                n.VotedFor = n.ID     // Проголосовать за себя
                n.mu.Unlock()

                // Запуск выборов
                n.Inbox <- Message{Type:"RequestVote"}
            } else {
                n.mu.Unlock()
            }

            timeout = time.Duration(rand.Intn(150)+150) * time.Millisecond
            electionTimer.Reset(timeout)

        }
    }
}


func (n *Node) sendHeartbeats(nodes map[int]*Node) {
	//логика отправки сообщений Heartbeats
 }

func (n *Node) handleRequestVote(msg Message, nodes map[int]*Node) {
	if n.CurrentTerm < msg.Term && n.VotedFor == -1{
		n.VotedFor = msg.FromID
		nodes[msg.FromID].Inbox <- Message{Type:"RequestVoteReply", FromID: n.ID, ToID: msg.FromID, VoteGranted: true}
	}

	
}

func (n *Node) handleRequestVoteReply(msg Message, nodes map[int]*Node) {
	if n.countVotes == 0{
		n.countVotes++
		for _, peer := n.Peers{
			peer.Inbox <- Message{Type:"RequestVote",Term: n.CurrentTerm, FromID: n.ID, ToID: peer.ID}
		}
	}else{
		if msg.VoteGranted{
			n.countVotes
		}

		if n.countVotes >= (len(nodes) - 1) /2 {
			n.countVotes = 0
			
		}
	}
}

func (n *Node) handleAppendEntries(msg Message, nodes map[int]*Node) {

	// Если полученный термин больше текущего, обновляем текущий термин и переходим в состояние Follower
	if msg.Term > n.CurrentTerm {
		n.CurrentTerm = msg.Term
		n.State = Follower
		n.VotedFor = -1 // Сбрасываем голос
	}

	// Если термин меньше текущего, игнорируем сообщение
	if msg.Term < n.CurrentTerm {
		return
	}

	// Устанавливаем узел как Follower
	n.State = Follower

	// Проверяем, соответствует ли индекс последней записи в логе
	if msg.LastLogIndex >= len(n.Log) || (msg.LastLogIndex >= 0 && n.Log[msg.LastLogIndex].Term != msg.LastLogTerm) {
		// Если нет, отправляем ответ с успехом = false
		nodes[msg.FromID].Inbox <- Message{Type: "AppendEntriesReply", Term: n.CurrentTerm, Success: false}
		return
	}

	// Если все проверки пройдены, обновляем лог
	n.Log = n.Log[:msg.LastLogIndex+1] // Обрезаем лог до нужного индекса
	n.Log = append(n.Log, msg.Entries...) // Добавляем новые записи

	// Обновляем индекс коммита
	if msg.LeaderCommit > n.CommitIndex {
		n.CommitIndex = min(msg.LeaderCommit, len(n.Log)-1)
	}

	// Отправляем ответ с успехом = true
	nodes[msg.FromID].Inbox <- Message{Type: "AppendEntriesReply", Term: n.CurrentTerm, Success: true}
}

func (n *Node) handleAppendEntriesReply(msg Message, nodes map[int]*Node) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// Если полученный термин больше текущего, обновляем текущий термин и переходим в состояние Follower
	if msg.Term > n.CurrentTerm {
		n.CurrentTerm = msg.Term
		n.State = Follower
		n.VotedFor = -1 // Сбрасываем голос
		return
	}

	// Если термин меньше текущего, игнорируем сообщение
	if msg.Term < n.CurrentTerm {
		return
	}

	// Обработка успешного ответа
	if msg.Success {
		// Увеличиваем индекс совпадения для узла
		n.MatchIndex[msg.FromID] = len(n.Log) - 1
		n.NextIndex[msg.FromID] = len(n.Log) // Обновляем индекс для следующего AppendEntries

		// Обновляем CommitIndex, если есть необходимость
		n.updateCommitIndex()
	} else {
		// Если неуспех, уменьшаем индекс следующего
		n.NextIndex[msg.FromID]--
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}