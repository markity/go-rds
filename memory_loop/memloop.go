package memoryloop

import (
	"go-rds/commands"
)

func memloop(mem *memLoop) {
	for {
		mem.mu.Lock()
		if len(mem.commands) != 0 {
			tobeProcceed := mem.commands
			mem.commands = make([]*commands.MemLoopCommandCover, 0)
			mem.mu.Unlock()

			for _, command := range tobeProcceed {
				switch command.Command.(type) {
				case
					// handshake hello is not required, proceeded by on_message
					*commands.SetInfoLibNameCommand,
					*commands.SetInfoLibVersionCommand,
					*commands.PingCommand:

					command.Loop.RunInLoop(func() {
						MemLoopCallback(command)
					})
				}
			}
		} else {
			mem.cond.Wait()
			mem.mu.Unlock()
		}
	}
}
