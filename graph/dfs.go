package graph

func DFS(g *Graph, start int) []int {
	if len(g.Adj) == 0 {
		return []int{}
	}
	visited := make(map[int]bool)
	visited[start] = true
	order := []int{}
	stack_slice := &Stack{}
	stack_slice.Push(start)
	for !stack_slice.IsEmpty() {
		u, ok := stack_slice.Pop()
		if !ok {
			// Если очередь пуста (и произошла ошибка Dequeue), выходим.
			break
		}
		order = append(order, u)
		for _, neighbor := range g.Adj[u] {
			if !visited[neighbor] {
				visited[neighbor] = true
				stack_slice.Push(neighbor)
			}
		}
	}

	return order
}
