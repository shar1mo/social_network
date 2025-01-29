package main

import (
	"fmt"
	"time"
	"wintersc/distributed"
)

func main() {
	// Создаём 5 узлов
	nodeIDs := []int{0, 1, 2, 3, 4}
	var nodes []*distributed.Node

	for _, id := range nodeIDs {
		nodes = append(nodes, distributed.NewNode(id))
	}

	// Формируем кольцо
	distributed.SetupRing(nodes)

	// Запускаем Ring-based Election с узла 0
	go nodes[0].StartRingElection()

	// Даем время выбрать лидера
	time.Sleep(3 * time.Second)

	// Имитируем сбой: узел 3 "падает"
	nodes[3].Alive = false
	fmt.Println("Node 3 отключен!")

	// Даем лидеру собрать данные
	for _, node := range nodes {
		if node.LeaderID == node.ID {
			go node.StartGlobalCollection()
		}
	}

	time.Sleep(3 * time.Second)

	// Проверяем, кто стал лидером
	for _, node := range nodes {
		fmt.Printf("Node %d: Мой лидер %d, у меня %d пользователей\n", node.ID, node.LeaderID, node.LocalData)
	}
}
