package graph

type Graph struct {
	Adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{Adj: make(map[int][]int)}
}

func (g *Graph) AddEdge(u, v int) {
	if _, exists := g.Adj[u]; !exists {
		g.Adj[u] = []int{}
	}
	if _, exists := g.Adj[v]; !exists {
		g.Adj[v] = []int{}
	}
	g.Adj[u] = append(g.Adj[u], v)
	g.Adj[v] = append(g.Adj[v], u)
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
