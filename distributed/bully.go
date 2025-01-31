package distributed

import (
	"fmt"
	"sync"
	"time"
)

type BullyNode struct {
	HasNeighbour bool
	ID           int
	LeaderID     int
	Alive        bool
	LocalData    int // Число пользователей
	Inbox        chan Msg
	BullyNodes   []*BullyNode
	Mutex        sync.Mutex
	IsListen     bool
}

type Msg struct {
	Kind string     // "ELECTION", "COORDINATOR", "COLLECT", "COLLECT_REPLY"
	From *BullyNode // ID отправителя
	Data int        // Локальные данные (для сбора)
}

func NewBullyNode(id int) *BullyNode {
	return &BullyNode{
		IsListen:     false,
		ID:           id,
		LeaderID:     -1,
		Alive:        true,
		Inbox:        make(chan Msg, 10),
		BullyNodes:   make([]*BullyNode, 0),
		HasNeighbour: true,
	}
}

func (n *BullyNode) Listen() {
	timer := time.NewTimer(3 * time.Second)
	for {
		select {
		case msg := <-n.Inbox:
			switch msg.Kind {
			case "ELECTION":
				n.HandleElection(msg)
			case "COORDINATOR":
				n.HandleCoordinator(msg)
			case "OK":
				n.HandleOk(msg)

			}
		case <-timer.C:
			n.StartBullyElection()
		}
	}
}

func (node *BullyNode) AddConnectBully(nodes []*BullyNode) {
	node.BullyNodes = append(node.BullyNodes, nodes...)

	if !node.IsListen {
		go node.Listen()
	}
}

func SendMessage(to *BullyNode, msg Msg) {

	to.Inbox <- msg

}

func (n *BullyNode) StartBullyElection() {
	for _, neighbour := range n.BullyNodes {
		if neighbour.ID > n.ID {
			SendMessage(neighbour, Msg{Kind: "ELECTION", From: n})
		}
	}

	time.Sleep(2 * time.Second)

	if n.HasNeighbour {
		for _, neighbour := range n.BullyNodes {
			SendMessage(neighbour, Msg{Kind: "COORDINATOR", From: n})
			n.Mutex.Lock()
			n.LeaderID = n.ID
			n.Mutex.Unlock()
		}
	}
}

func (n *BullyNode) HandleOk(msg Msg) {
	fmt.Printf("%d: Пришел OK от %d\n", n.ID, msg.From.ID)
	if !n.Alive {
		n.Mutex.Lock()
		n.HasNeighbour = false
		n.Mutex.Unlock()
	}
	if n.ID < msg.From.ID {
		n.Mutex.Lock()
		n.HasNeighbour = false
		n.Mutex.Unlock()
		fmt.Println(n.HasNeighbour)
	} else {
		n.Mutex.Lock()
		n.HasNeighbour = true
		n.Mutex.Unlock()
		fmt.Println(n.HasNeighbour)
	}

}

func (n *BullyNode) HandleElection(msg Msg) {
	fmt.Printf("%d: Пришел ELECTION от %d\n", n.ID, msg.From.ID)
	if n.Alive {
		SendMessage(msg.From, Msg{Kind: "OK", From: n})

		for _, neighbour := range n.BullyNodes {
			if neighbour.ID > n.ID {
				SendMessage(neighbour, Msg{Kind: "ELECTION", From: n})
			}
		}
	} else {
		n.HasNeighbour = false
	}
}

func (n *BullyNode) HandleCoordinator(msg Msg) {
	fmt.Printf("%d: Пришел COORDINATOR от %d\n", n.ID, msg.From.ID)
	n.Mutex.Lock()
	n.LeaderID = msg.From.ID
	n.Mutex.Unlock()
}
