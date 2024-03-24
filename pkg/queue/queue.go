package queue

// 定义一个简单的栈结构
type Queue []interface{}

func (s *Queue) PushBack(v interface{}) {
	*s = append(*s, v)
}

func (s *Queue) PopBack() interface{} {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *Queue) PopFront() interface{} {
	v := (*s)[0]
	*s = (*s)[1:]
	return v
}

func (s *Queue) Size() int {
	return len(*s)
}
