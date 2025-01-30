package main

import (
	"fmt"
	"time"
	"wintersc/distributed"
)

func main() {
	// Создаём 4 узла
	nodeIDs := []int{0, 1, 2, 3}
	var nodes []*distributed.BullyNode

	for _, id := range nodeIDs {
		nodes = append(nodes, distributed.NewBullyNode(id))
	}

	// Формируем связи Bully
	nodes[0].AddConnectBully([]*distributed.BullyNode{nodes[1], nodes[2], nodes[3]})
	nodes[1].AddConnectBully([]*distributed.BullyNode{nodes[0], nodes[2], nodes[3]})
	nodes[2].AddConnectBully([]*distributed.BullyNode{nodes[0], nodes[1], nodes[3]})
	nodes[3].AddConnectBully([]*distributed.BullyNode{nodes[0], nodes[1], nodes[2]})

	// Запускаем Listen() для каждого узла в отдельной горутине

	// Запускаем Bully с узла 0
	nodes[2].Alive = false
	nodes[3].Alive = false

	go nodes[0].StartBullyElection()
	// Даем время выбрать лидера
	time.Sleep(15 * time.Second)

	// Проверяем, кто стал лидером
	for _, node := range nodes {
		fmt.Printf("Node %d: Мой лидер %d\n", node.ID, node.LeaderID)
	}
}
