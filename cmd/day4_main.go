package main

import (
	"time"
	"wintersc/distributed"
)

func main() {
	// Создаём 3 узла Raft
	node1 := distributed.NewRaftNode(1, nil)
	node2 := distributed.NewRaftNode(2, nil)
	node3 := distributed.NewRaftNode(3, nil)

	// Связываем узлы
	node1.Nodes = []*distributed.RaftNode{node2, node3}
	node2.Nodes = []*distributed.RaftNode{node1, node3}
	node3.Nodes = []*distributed.RaftNode{node1, node2}

	// Запускаем горутины для каждого узла
	go node1.drun()
	go node2.run()
	go node3.run()

	// Симулируем клиент, который отправляет команды лидеру
	time.Sleep(2 * time.Second)
	node1.distributedsendMessage(node1.LeaderID, "AppendEntries", node1.Term, "commandX")

	// Падение лидера (например, узел 1)
	time.Sleep(3 * time.Second)
	node1.distributedfailLeader()

	// Ждем, пока новый лидер не станет лидером
	time.Sleep(2 * time.Second)

	// Снова отправляем команду
	node2.distributed.sendMessage(node2.LeaderID, "AppendEntries", node2.Term, "commandY")
}
