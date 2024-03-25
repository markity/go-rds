package memoryloop

import (
	"go-rds/commands"
	datastructure "go-rds/data_structure"
	"sync"
	"sync/atomic"
	"time"
)

type memLoop struct {
	bigKV   map[string]*datastructure.RdsObject
	started atomic.Int64

	commands []*commands.MemLoopCommandCover
	mu       sync.Mutex
	cond     *sync.Cond
}

// pay attention to the ttl
func (mem *memLoop) GetRdsObj(key string) (*datastructure.RdsObject, time.Time) {
	now := time.Now()
	rdsObj, ok := mem.bigKV[key]
	if !ok {
		return nil, now
	}

	if rdsObj.TTL != nil && now.After(*rdsObj.TTL) {
		delete(mem.bigKV, key)
		return nil, now
	}

	return mem.bigKV[key], now
}

func (mem *memLoop) DelRdsObj(key string) {
	delete(mem.bigKV, key)
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
