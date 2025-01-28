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
	var edges []Edge
	// Перебираем каждую вершину
	for u, neighbors := range g.Adj {
		// Для каждого соседа вершины u
		for _, v := range neighbors {
			// Если u < v, добавляем ребро (u, v), иначе пропускаем
			if u < v {
				// Мы добавляем информацию о весе ребра. Допустим, что веса уже хранятся в g.Edge.
				// Но если вес не был добавлен, то его можно добавить вручную при формировании ребра
				var weight int // Вес можно добавить, если информация о весах имеется
				// Пример: поиск веса ребра между u и v
				for _, edge := range g.Edge {
					if (edge.U == u && edge.V == v) || (edge.U == v && edge.V == u) {
						weight = edge.W
						break
					}
				}
				edges = append(edges, Edge{U: u, V: v, W: weight})
			}
		}
	}
	return edges
}
