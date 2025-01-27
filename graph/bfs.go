package graph

func BFS(g *Graph, start int) []int {
	visited := make(map[int]bool)
	visited[start] = true
	order := []int{}
	queue_slice := &Queue{}
	queue_slice.Enqueue(start)
	for !queue_slice.IsEmpty() {
		u, ok := queue_slice.Dequeue()
		if !ok {
			// Если очередь пуста (и произошла ошибка Dequeue), выходим.
			break
		}
		order = append(order, u)
		for _, neighbor := range g.Adj[u] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue_slice.Enqueue(neighbor)
			}
		}
	}

	return order
}
