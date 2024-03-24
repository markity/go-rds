package memoryloop

import (
	"go-rds/commands"
	datastructure "go-rds/data_structure"
	"sync"
	"sync/atomic"
)

type memLoop struct {
	bigKV   map[string]*datastructure.RdsObject
	started atomic.Int64

	commands []*commands.MemLoopCommandCover
	mu       sync.Mutex
	cond     *sync.Cond
}

func (mem *memLoop) QueueCommand(command *commands.MemLoopCommandCover) {
	mem.mu.Lock()
	mem.commands = append(mem.commands, command)
	mem.cond.Signal()
	mem.mu.Unlock()
}

func (mem *memLoop) StartLoop() {
	ok := mem.started.CompareAndSwap(0, 1)
	if !ok {
		panic("already in loop")
	}

	go func() {
		memloop(mem)
	}()
}

func NewMemLoop() MemLoop {
	memLoop := &memLoop{
		bigKV: make(map[string]*datastructure.RdsObject),
	}
	memLoop.cond = sync.NewCond(&memLoop.mu)
	return memLoop
}
