package graph

type Stack struct {
	Data []int
}

func (s *Stack) Push(x int) {
	s.Data = append(s.Data, x)
}

func (s *Stack) Pop() (int, bool) {
	if s.IsEmpty() {
		return 0, false
	}
	elem := s.Data[len(s.Data)-1]
	s.Data = s.Data[:len(s.Data)-1]
	return elem, true
}

func (s *Stack) IsEmpty() bool {
	return len(s.Data) == 0
}
