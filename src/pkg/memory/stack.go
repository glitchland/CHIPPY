package memory

type Stack struct {
	stack []uint16
	Top uint8
	SP uint8
}

func (s *Stack) Init() {
}

func (s *Stack) Push(v uint16) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() uint16 {
	if s.SP < 0 {
		panic("Stack underflow triggered")
	}
	v := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return v
}
