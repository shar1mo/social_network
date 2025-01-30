package distributed

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RingNode struct {
	ID        int
	NextID    int
	LeaderID  int
	Alive     bool
	LocalData int // Число пользователей
	Inbox     chan Message
	RingNodes map[int]*RingNode
	Mutex     sync.Mutex
}

type Message struct {
	Kind   string // "ELECTION", "COORDINATOR", "COLLECT", "COLLECT_REPLY"
	IDs    []int  // Список ID узлов (для выборов)
	FromID int    // ID отправителя
	Data   int    // Локальные данные (для сбора)
}

func NewRingNode(id int) *RingNode {
	return &RingNode{
		ID:        id,
		NextID:    -1,
		LeaderID:  -1,
		Alive:     true,
		LocalData: rand.Intn(50) + 50, // Рандомное количество пользователей (50-100)
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
			case "COLLECT":
				n.HandleCollect(msg)
			case "COLLECT_REPLY":
				n.HandleCollectReply(msg)
			}
			// Вы можете добавить обработку других каналов, если это необходимо
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

// Лидер запускает сбор данных
func (n *RingNode) StartGlobalCollection() {
	fmt.Printf("RingNode %d (Лидер): Начинаю сбор данных\n", n.ID)
	n.Mutex.Lock()
	expectedReplies := len(n.RingNodes) - 1
	n.Mutex.Unlock()

	received := 0
	sum := 0

	for _, RingNode := range n.RingNodes {
		if RingNode.ID != n.ID && RingNode.Alive {
			n.SendMessage(RingNode.ID, Message{Kind: "COLLECT", FromID: n.ID})
		}
	}

	timeout := time.After(2 * time.Second)
	for received < expectedReplies {
		select {
		case msg := <-n.Inbox:
			if msg.Kind == "COLLECT_REPLY" {
				fmt.Printf("RingNode %d: Получил данные от RingNode %d: %d\n", n.ID, msg.FromID, msg.Data)
				sum += msg.Data
				received++
			}
		case <-timeout:
			fmt.Println("Лидер: Превышено время ожидания ответа")
			break
		}
	}

	fmt.Printf("RingNode %d (Лидер): Общая сумма пользователей: %d\n", n.ID, sum)
}

// Обработка COLLECT-запроса
func (n *RingNode) HandleCollect(msg Message) {
	if n.Alive {
		fmt.Printf("RingNode %d: Отправляю свои данные лидеру %d\n", n.ID, msg.FromID)
		n.SendMessage(msg.FromID, Message{Kind: "COLLECT_REPLY", FromID: n.ID, Data: n.LocalData})
	}
}

// Обработка COLLECT_REPLY лидером
func (n *RingNode) HandleCollectReply(msg Message) {
	fmt.Printf("RingNode %d (Лидер): Получил ответ от %d: %d\n", n.ID, msg.FromID, msg.Data)
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
