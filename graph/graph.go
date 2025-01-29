package graph

type Graph struct {
	Adj  map[int][]int
	Edge []Edge
}

type Edge struct {
	U, V, W int
}

func NewGraph() *Graph {
	return &Graph{
		Adj:  make(map[int][]int),
		Edge: []Edge{},
	}
}

func (g *Graph) AddEdge(u, v, w int) {
	if _, exists := g.Adj[u]; !exists {
		g.Adj[u] = []int{}
	}
	if _, exists := g.Adj[v]; !exists {
		g.Adj[v] = []int{}
	}
	g.Adj[u] = append(g.Adj[u], v)
	g.Adj[v] = append(g.Adj[v], u)

	g.Edge = append(g.Edge, Edge{U: u, V: v, W: w})
}

func HasEdge(g *Graph, u, v int) bool {
	if neighbors, exists := g.Adj[u]; exists {
		for _, neighbor := range neighbors {
			if neighbor == v {
				return true
			}
		}
	}
	return false
}

func ConnectedComponents(g *Graph) (count int, comp map[int]int) {
	visited := make(map[int]bool) // Для отслеживания посещённых узлов
	comp = make(map[int]int)      // Для хранения компонент связности
	count = 0                     // Счётчик компонент связности

	// Перебираем все узлы графа
	for key := range g.Adj {
		if !visited[key] {
			count++ // Новая компонента связности
			// Получаем все узлы компоненты с помощью DFS
			component := DFS(g, key)
			// Обрабатываем все узлы из компоненты
			for _, value := range component {
				visited[value] = true // Помечаем узел как посещённый
				comp[value] = count   // Присваиваем номер компоненты
			}
		}
	}

	return count, comp
}

func (g *Graph) GetAllEdges() []Edge {
	// Создаём карту для быстрого поиска веса ребра
	weights := make(map[[2]int]int)
	for _, edge := range g.Edge {
		weights[[2]int{edge.U, edge.V}] = edge.W
		weights[[2]int{edge.V, edge.U}] = edge.W // Для неориентированного графа
	}

	var edges []Edge
	// Перебираем каждую вершину
	for u, neighbors := range g.Adj {
		for _, v := range neighbors {
			// Избегаем дублирования рёбер, добавляя только (u, v) где u < v
			if u < v {
				weight := weights[[2]int{u, v}]
				edges = append(edges, Edge{U: u, V: v, W: weight})
			}
		}
	}
	return edges
}

func (g *Graph) GetNeighbors(u int) []struct{ V, W int } {
	var neighbors []struct{ V, W int }

	for _, v := range g.Adj[u] {
		for _, edge := range g.Edge {
			if (edge.U == u && edge.V == v) || (edge.U == v && edge.V == u) {
				neighbors = append(neighbors, struct{ V, W int }{V: v, W: edge.W})
				break
			}
		}
	}
	return neighbors
}
