package distributed

import (
	"fmt"
	"sync"
)

type RingNode struct {
	ID        int
	NextID    int
	LeaderID  int
	Alive     bool
	Inbox     chan Message
	RingNodes map[int]*RingNode
	Mutex     sync.Mutex
}

type Message struct {
	Kind   string
	IDs    []int // Список ID узлов (для выборов)
	FromID int   // ID отправителя
	Data   int   // Локальные данные (для сбора)
}

func NewRingNode(id int) *RingNode {
	return &RingNode{
		ID:        id,
		NextID:    -1,
		LeaderID:  -1,
		Alive:     true,
		Inbox:     make(chan Message, 10),
		RingNodes: make(map[int]*RingNode),
	}
}

// Устанавливаем кольцевые связи
func SetupRing(RingNodes []*RingNode) {
	for i, RingNode := range RingNodes {
		nextIndex := (i + 1) % len(RingNodes)
		RingNode.NextID = RingNodes[nextIndex].ID
		for _, n := range RingNodes {
			RingNode.RingNodes[n.ID] = n
		}
		go RingNode.Listen()
	}
}

func (n *RingNode) Listen() {
	for {
		select {
		case msg := <-n.Inbox:
			switch msg.Kind {
			case "ELECTION":
				n.HandleElection(msg)
			case "COORDINATOR":
				n.HandleCoordinator(msg)
			}
		}
	}
}

func (n *RingNode) StartRingElection() {
	fmt.Printf("RingNode %d: Начинаю кольцевые выборы\n", n.ID)
	n.SendMessage(n.NextID, Message{Kind: "ELECTION", IDs: []int{n.ID}, FromID: n.ID})
}

// Обработка ELECTION в кольце
func (n *RingNode) HandleElection(msg Message) {
	if !n.Alive {
		// Если узел неактивен, пропускаем его
		n.SendMessage(n.NextID, msg)
		return
	}

	maxID := msg.IDs[0]
	if n.ID > maxID {
		msg.IDs = []int{n.ID}
	}

	if msg.FromID == n.ID {
		fmt.Printf("RingNode %d: Выбран лидер %d, рассылаю COORDINATOR\n", n.ID, maxID)
		n.Broadcast(Message{Kind: "COORDINATOR", FromID: maxID})
	} else {
		n.SendMessage(n.NextID, msg)
	}
}

// Обработка COORDINATOR
func (n *RingNode) HandleCoordinator(msg Message) {
	fmt.Printf("RingNode %d: Принял нового лидера %d\n", n.ID, msg.FromID)
	n.Mutex.Lock()
	n.LeaderID = msg.FromID
	n.Mutex.Unlock()

	if msg.FromID != n.ID {
		n.SendMessage(n.NextID, msg)
	}
}

// Отправка сообщения узлу
func (n *RingNode) SendMessage(to int, msg Message) {
	if RingNode, ok := n.RingNodes[to]; ok {
		if RingNode.Alive { // Проверяем, активен ли узел
			RingNode.Inbox <- msg
		} else {
			fmt.Printf("RingNode %d: Узел %d не активен, пропускаю сообщение\n", n.ID, to)
			n.SendMessage(RingNode.NextID, msg) // Пропускаем неактивный узел
		}
	}
}

func (n *RingNode) Broadcast(msg Message) {
	for _, RingNode := range n.RingNodes {
		RingNode.Inbox <- msg
	}
}
