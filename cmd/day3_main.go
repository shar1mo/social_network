package main

import (
	"fmt"
	"time"
	"wintersc/distributed"
)

func main() {
	// Создаём 4 узла
	nodeIDs := []int{0, 1, 2, 3, 7, 9}
	var nodes []*distributed.RingNode

	for _, id := range nodeIDs {
		nodes = append(nodes, distributed.NewRingNode(id))
	}

	// Формируем связи Bully
	distributed.SetupRing(nodes)

	// Запускаем Listen() для каждого узла в отдельной горутине

	// Запускаем Bully с узла 0
	nodes[5].Alive = false
	go nodes[0].StartRingElection()
	// Даем время выбрать лидера
	time.Sleep(6 * time.Second)

	// Проверяем, кто стал лидером
	for _, node := range nodes {
		fmt.Printf("Node %d: Мой лидер %d\n", node.ID, node.LeaderID)
	}
}
