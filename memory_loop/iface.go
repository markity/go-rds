package memoryloop

import "go-rds/commands"

type MemLoop interface {
	QueueCommand(command *commands.MemLoopCommandCover)
	StartLoop()
}
