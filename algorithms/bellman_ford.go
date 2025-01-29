package algorithms

import (
	"math"
	"wintersc/graph"
)

func BellmanFord(g *graph.Graph, start int) (map[int]int, map[int]int, bool) {
	// Инициализация расстояний и предков
	dist := make(map[int]int)
	prev := make(map[int]int)

	// Устанавливаем все расстояния как "бесконечность"
	for vertex := range g.Adj {
		dist[vertex] = math.MaxInt32
		prev[vertex] = -1
	}
	dist[start] = 0

	// Основной цикл алгоритма (проходим |V| - 1 раз)
	for i := 1; i < len(g.Adj); i++ {
		for _, edge := range g.Edge {
			// Если найден более короткий путь через ребро
			if dist[edge.U] != math.MaxInt32 && dist[edge.U]+edge.W < dist[edge.V] {
				dist[edge.V] = dist[edge.U] + edge.W
				prev[edge.V] = edge.U
			}
		}
	}

	// Проверка на наличие отрицательных циклов
	negativeCycle := false
	for _, edge := range g.Edge {
		if dist[edge.U] != math.MaxInt32 && dist[edge.U]+edge.W < dist[edge.V] {
			negativeCycle = true
			break
		}
	}

	return dist, prev, negativeCycle
}
