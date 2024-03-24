package global

import memoryloop "go-rds/memory_loop"

var MemLoop memoryloop.MemLoop

func init() {
	MemLoop = memoryloop.NewMemLoop()
	MemLoop.StartLoop()
}
