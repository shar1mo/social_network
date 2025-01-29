package graph

type Item struct {
	Vertex, Dist int
}

type PriorityQueue struct {
	Data []Item
	Item Item
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{}
}

func (q *PriorityQueue) Push(Vertex, dist int) {
	item := Item{Vertex: Vertex, Dist: dist}
	q.Data = append(q.Data, item)
	q.HeapfiUp()
}

func (q *PriorityQueue) HeapfiUp() {
	i := len(q.Data) - 1

	for i > 0 {
		parentIndex := (i - 1) / 2

		if q.Data[parentIndex].Vertex > q.Data[i].Vertex {
			q.Data[parentIndex], q.Data[i] = q.Data[i], q.Data[parentIndex]
			i = parentIndex
		} else {
			break
		}
	}
}

func (q *PriorityQueue) Pop() (Item, bool) {
	if len(q.Data) == 0 {
		return Item{}, false
	}
	min_el := q.Data[0]
	q.HeapfiDown()
	return min_el, true
}

func (q *PriorityQueue) HeapfiDown() {
	if len(q.Data) == 0 {
		panic("Pop from an empty priority queue")
	}
	q.Data[0] = q.Data[len(q.Data)-1]
	q.Data = q.Data[:len(q.Data)-1]

	i := 0
	for {
		leftIndex := 2*i + 1
		rightIndex := 2*i + 2
		smallestIndex := i

		if leftIndex < len(q.Data) && q.Data[leftIndex].Vertex < q.Data[smallestIndex].Vertex {
			smallestIndex = leftIndex
		}

		if rightIndex < len(q.Data) && q.Data[rightIndex].Vertex < q.Data[smallestIndex].Vertex {
			smallestIndex = rightIndex
		}

		if smallestIndex == i {
			break
		}

		q.Data[i], q.Data[smallestIndex] = q.Data[smallestIndex], q.Data[i]
		i = smallestIndex
	}

}
