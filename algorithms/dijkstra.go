package algorithms

import (
	"math"
	"wintersc/graph"
)

func Dijkstra(g *graph.Graph, start int) (map[int]int, map[int]int) {
	// Создаём таблицу расстояний, изначально все - "бесконечность"
	distances := make(map[int]int)
	for _, edge := range g.Edge {
		distances[edge.U] = math.MaxInt32
		distances[edge.V] = math.MaxInt32
	}
	distances[start] = 0 // Начальная вершина

	// Массив для восстановления пути
	prev := make(map[int]int)

	// Создаём приоритетную очередь
	pq := graph.NewPriorityQueue()
	pq.Push(start, 0)

	// Основной цикл алгоритма
	for len(pq.Data) > 0 {
		// Берём вершину с минимальным расстоянием
		current, ok := pq.Pop()
		if !ok {
			break
		}

		currentNode := current.Vertex
		currentDistance := current.Dist

		// Если найденный путь длиннее, чем уже известный, пропускаем
		if currentDistance > distances[currentNode] {
			continue
		}

		// Обновляем расстояния до соседей
		for _, edge := range g.GetNeighbors(currentNode) {
			newDistance := currentDistance + edge.W
			if newDistance < distances[edge.V] {
				distances[edge.V] = newDistance
				prev[edge.V] = currentNode
				pq.Push(edge.V, newDistance)
			}
		}
	}

	return distances, prev
}
