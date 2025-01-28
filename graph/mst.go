package graph

func BoruvkaMST(n int, edges []Edge) (mst []Edge, totalWeight int) {
	ds := NewDisjointSet(n)
	mst = []Edge{}
	totalWeight = 0

	numComponents := n

	for numComponents > 1 {
		minEdges := make([]Edge, n)

		for i := range minEdges {
			minEdges[i] = Edge{-1, -1, int(^uint(0) >> 1)}
		}

		for _, edge := range edges {
			u, v, w := edge.U, edge.V, edge.W

			compU := ds.Find(u)
			compV := ds.Find(v)

			if compU != compV {
				if w < minEdges[compU].W {
					minEdges[compU] = edge
				}
				if w < minEdges[compV].W {
					minEdges[compV] = edge
				}
			}
		}

		for _, edge := range minEdges {
			if edge.U == -1 && edge.V == -1 {
				continue
			}

			u, v, w := edge.U, edge.V, edge.W

			if ds.Find(u) != ds.Find(v) {
				ds.Union(u, v)
				mst = append(mst, edge)
				totalWeight += w
				numComponents--
			}
		}
	}
	return mst, totalWeight
}

func MergeSort(edges []Edge) []Edge {
	if len(edges) <= 1 {
		return edges
	}
	mid := len(edges) / 2
	left := MergeSort(edges[:mid])
	right := MergeSort(edges[mid:])

	return Merge(left, right)
}

func Merge(left, right []Edge) []Edge {
	var result []Edge
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].W <= right[j].W {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	for i < len(left) {
		result = append(result, left[i])
		i++
	}

	for j < len(right) {
		result = append(result, right[j])
		j++
	}
	return result
}
