package graph

type DisjointSet struct {
	Parent []int
	Rank   []int
}

func NewDisjointSet(n int) *DisjointSet {
	parent := make([]int, n)
	rank := make([]int, n)

	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 0
	}
	return &DisjointSet{
		Parent: parent,
		Rank:   rank,
	}
}

func (ds *DisjointSet) Find(x int) int {
	if ds.Parent[x] != x {
		ds.Parent[x] = ds.Find(ds.Parent[x])
	}
	return ds.Parent[x]
}

func (ds *DisjointSet) Union(x, y int) bool {
	rootX := ds.Find(x)
	rootY := ds.Find(y)

	if rootX == rootY {
		return false
	}

	if ds.Rank[rootX] > ds.Rank[rootY] {
		ds.Parent[rootY] = rootX
	} else if ds.Rank[rootX] < ds.Rank[rootY] {
		ds.Parent[rootX] = rootY
	} else {
		ds.Rank[rootY] = rootX
		ds.Rank[rootX]++
	}
	return true
}
