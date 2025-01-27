package graph

type Queue struct {
	Data []int
}

func (s *Queue) Enqueue(x int) {
	s.Data = append(s.Data, x)
}

func (q *Queue) Dequeue() (int, bool) {
	if q.IsEmpty() {
		return 0, false
	}
	tmp := q.Data[0]
	q.Data = q.Data[1:]
	return tmp, true
}

func (q *Queue) IsEmpty() bool {
	return len(q.Data) == 0
}
