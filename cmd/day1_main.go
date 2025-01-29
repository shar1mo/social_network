package main

import (
	"fmt"
	"wintersc/graph"
)

func main() {
	graph1 := graph.NewGraph()

	graph1.AddEdge(1, 2, 4)
	graph1.AddEdge(3, 1, 7)
	graph1.AddEdge(2, 4, 1)
	graph1.AddEdge(1, 5, 3)
	graph1.AddEdge(6, 7, 2)

	fmt.Println(graph.HasEdge(graph1, 1, 5))

	stack := &graph.Stack{}

	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	fmt.Println(stack)

	stack.Pop()
	fmt.Println(stack)

	queue := &graph.Queue{}
	queue.Enqueue(10)
	queue.Enqueue(20)
	queue.Enqueue(30)

	fmt.Println(queue)

	queue.Dequeue()

	fmt.Println(queue)

	fmt.Printf("%d\n\n", graph1.Adj)

	fmt.Println(graph.BFS(graph1, 1))

	fmt.Println(graph.DFS(graph1, 1))

	fmt.Println(graph.ConnectedComponents(graph1))
}
