package main

import (
	"fmt"
	"wintersc/graph"
)

func main() {
	graph1 := graph.NewGraph()

	graph1.AddEdge(1, 2)
	graph1.AddEdge(1, 3)
	graph1.AddEdge(1, 4)
	graph1.AddEdge(2, 4)
	graph1.AddEdge(3, 8)
	graph1.AddEdge(2, 7)

	graph1.AddEdge(11, 12)

	//fmt.Printf("%d\n\n", graph.Adj)

	fmt.Println(graph.BFS(graph1, 1))

	fmt.Println(graph.ConnectedComponents(graph1))
}
