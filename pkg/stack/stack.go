package stack

// 定义一个简单的栈结构
type Stack []interface{}

func (s *Stack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *Stack) Pop() interface{} {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *Stack) Peek() interface{} {
	return (*s)[len(*s)-1]
}

func (s *Stack) Size() int {
	return len(*s)
}
