package distributed

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	ID        int
	NextID    int
	LeaderID  int
	Alive     bool
	LocalData int // Число пользователей
	Inbox     chan Message
	Nodes     map[int]*Node
	Mutex     sync.Mutex
}

type Message struct {
	Kind   string // "ELECTION", "COORDINATOR", "COLLECT", "COLLECT_REPLY"
	IDs    []int  // Список ID узлов (для выборов)
	FromID int    // ID отправителя
	Data   int    // Локальные данные (для сбора)
}

func NewNode(id int) *Node {
	return &Node{
		ID:        id,
		NextID:    -1,
		LeaderID:  -1,
		Alive:     true,
		LocalData: rand.Intn(50) + 50, // Рандомное количество пользователей (50-100)
		Inbox:     make(chan Message, 10),
		Nodes:     make(map[int]*Node),
	}
}

// Устанавливаем кольцевые связи
func SetupRing(nodes []*Node) {
	for i, node := range nodes {
		nextIndex := (i + 1) % len(nodes)
		node.NextID = nodes[nextIndex].ID
		for _, n := range nodes {
			node.Nodes[n.ID] = n
		}
		go node.Listen()
	}
}

// Узел слушает входящие сообщения
func (n *Node) Listen() {
	for msg := range n.Inbox {
		switch msg.Kind {
		case "ELECTION":
			n.HandleElection(msg)
		case "COORDINATOR":
			n.HandleCoordinator(msg)
		case "COLLECT":
			n.HandleCollect(msg)
		case "COLLECT_REPLY":
			n.HandleCollectReply(msg)
		}
	}
}

// Запуск Bully-алгоритма
func (n *Node) Bully() {
	fmt.Printf("Node %d: Запускаю Bully-выборы\n", n.ID)
	hasHigherNode := false

	for _, node := range n.Nodes {
		if node.ID > n.ID && node.Alive {
			n.SendMessage(node.ID, Message{Kind: "ELECTION", FromID: n.ID})
			hasHigherNode = true
		}
	}

	if !hasHigherNode {
		fmt.Printf("Node %d: Никто не ответил, становлюсь лидером\n", n.ID)
		n.Mutex.Lock()
		n.LeaderID = n.ID
		n.Mutex.Unlock()
		n.Broadcast(Message{Kind: "COORDINATOR", FromID: n.ID})
	}
}

// Запуск Ring-based Election
func (n *Node) StartRingElection() {
	fmt.Printf("Node %d: Начинаю кольцевые выборы\n", n.ID)
	n.SendMessage(n.NextID, Message{Kind: "ELECTION", IDs: []int{n.ID}, FromID: n.ID})
}

// Обработка ELECTION в кольце
func (n *Node) HandleElection(msg Message) {
	maxID := msg.IDs[0]
	if n.ID > maxID {
		msg.IDs = []int{n.ID}
	}

	if msg.FromID == n.ID {
		fmt.Printf("Node %d: Выбран лидер %d, рассылаю COORDINATOR\n", n.ID, maxID)
		n.Broadcast(Message{Kind: "COORDINATOR", FromID: maxID})
	} else {
		n.SendMessage(n.NextID, msg)
	}
}

// Обработка COORDINATOR
func (n *Node) HandleCoordinator(msg Message) {
	fmt.Printf("Node %d: Принял нового лидера %d\n", n.ID, msg.FromID)
	n.Mutex.Lock()
	n.LeaderID = msg.FromID
	n.Mutex.Unlock()

	if msg.FromID != n.ID {
		n.SendMessage(n.NextID, msg)
	}
}

// Лидер запускает сбор данных
func (n *Node) StartGlobalCollection() {
	fmt.Printf("Node %d (Лидер): Начинаю сбор данных\n", n.ID)
	n.Mutex.Lock()
	expectedReplies := len(n.Nodes) - 1
	n.Mutex.Unlock()

	received := 0
	sum := 0

	for _, node := range n.Nodes {
		if node.ID != n.ID && node.Alive {
			n.SendMessage(node.ID, Message{Kind: "COLLECT", FromID: n.ID})
		}
	}

	timeout := time.After(2 * time.Second)
	for received < expectedReplies {
		select {
		case msg := <-n.Inbox:
			if msg.Kind == "COLLECT_REPLY" {
				fmt.Printf("Node %d: Получил данные от Node %d: %d\n", n.ID, msg.FromID, msg.Data)
				sum += msg.Data
				received++
			}
		case <-timeout:
			fmt.Println("Лидер: Превышено время ожидания ответа")
			break
		}
	}

	fmt.Printf("Node %d (Лидер): Общая сумма пользователей: %d\n", n.ID, sum)
}

// Обработка COLLECT-запроса
func (n *Node) HandleCollect(msg Message) {
	if n.Alive {
		fmt.Printf("Node %d: Отправляю свои данные лидеру %d\n", n.ID, msg.FromID)
		n.SendMessage(msg.FromID, Message{Kind: "COLLECT_REPLY", FromID: n.ID, Data: n.LocalData})
	}
}

// Обработка COLLECT_REPLY лидером
func (n *Node) HandleCollectReply(msg Message) {
	fmt.Printf("Node %d (Лидер): Получил ответ от %d: %d\n", n.ID, msg.FromID, msg.Data)
}

// Отправка сообщения узлу
func (n *Node) SendMessage(to int, msg Message) {
	if node, ok := n.Nodes[to]; ok {
		node.Inbox <- msg
	}
}

// Рассылка всем узлам
func (n *Node) Broadcast(msg Message) {
	for _, node := range n.Nodes {
		node.Inbox <- msg
	}
}
