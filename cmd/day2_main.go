package main

import (
	"fmt"
	"wintersc/algorithms"
	"wintersc/graph"
)

func main() {
	// Инициализация DisjointSet
	fmt.Println("Инициализация DisjointSet:")
	n := 6
	ds := graph.NewDisjointSet(n)

	ds.Union(0, 1)
	ds.Union(1, 2)
	ds.Union(3, 4)

	fmt.Println("Результаты Find для всех вершин:")
	for i := 0; i < n; i++ {
		fmt.Printf("Find(%d): %d\n", i, ds.Find(i))
	}

	fmt.Println("\nUnion(2, 3):")
	ds.Union(2, 3)

	fmt.Println("\nМассивы parent и rank:")
	fmt.Printf("parent: %+v\n", ds.Parent)
	fmt.Printf("rank: %+v\n", ds.Rank)

	// Создаём граф
	fmt.Println("\nСоздаём граф:")
	g := graph.NewGraph()
	g.AddEdge(0, 1, 4)
	g.AddEdge(1, 2, 3)
	g.AddEdge(2, 3, 2)
	g.AddEdge(3, 4, 6)
	g.AddEdge(4, 5, 5)
	g.AddEdge(0, 5, 7)

	edges := g.GetAllEdges()
	fmt.Printf("Список рёбер: %+v\n", edges)

	neighbors := g.GetNeighbors(0)
	fmt.Println("Neighbors of 0:")
	for _, n := range neighbors {
		fmt.Printf("vertex: %d, weight: %d\n", n.V, n.W)
	}

	// MST
	fmt.Println("\nВычисляем MST:")
	mst, totalWeight := graph.BoruvkaMST(len(g.Adj), edges)
	fmt.Printf("Минимальное остовное дерево (MST): %+v\n", mst)
	fmt.Printf("Общий вес MST: %d\n", totalWeight)

	// Dijkstra's Algorithm
	fmt.Println("\nТест алгоритма Dijkstra:")
	start := 0
	distances, prev := algorithms.Dijkstra(g, start)
	fmt.Printf("Расстояния от вершины %d: %+v\n", start, distances)
	fmt.Printf("Предыдущие вершины: %+v\n", prev)

	// Bellman-Ford Algorithm
	fmt.Println("\nТест алгоритма Bellman-Ford:")
	distancesBF, prevBF, negativeCycle := algorithms.BellmanFord(g, start)
	if negativeCycle {
		fmt.Println("Обнаружен отрицательный цикл!")
	} else {
		fmt.Printf("Расстояния от вершины %d: %+v\n", start, distancesBF)
		fmt.Printf("Предыдущие вершины: %+v\n", prevBF)
	}
}
