package commands

import (
	eventloop "github.com/markity/go-reactor/pkg/event_loop"
)

type MemLoopCommandCover struct {
	Loop         eventloop.EventLoop
	ConnectionID int64
	Command      interface{}
}
