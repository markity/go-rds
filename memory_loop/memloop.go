package memoryloop

import (
	"go-rds/commands"
	datastructure "go-rds/data_structure"
	"go-rds/pkg/resp"
	"log"
	"time"
)

func memloop(mem *memLoop) {
	for {
		mem.mu.Lock()
		if len(mem.commands) != 0 {
			tobeProcceed := mem.commands
			mem.commands = make([]*commands.MemLoopCommandCover, 0)
			mem.mu.Unlock()

			for _, command := range tobeProcceed {
				switch innerCmd := command.Command.(type) {
				case *commands.SetInfoLibNameCommand, *commands.SetInfoLibVersionCommand:
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToBulkString("OK"))
					})
				case *commands.PingCommand:
					var msg string = "PONG"
					if innerCmd.Message != nil {
						msg = *innerCmd.Message
					}
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToBulkString(msg))
					})
				case *commands.HelloCommand:
					if innerCmd.Proto != "3" {
						log.Printf("hand shake failed\n")
						be, err := resp.ToBulkError("ERR", "proto number is not supoorted")
						if err != nil {
							panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, be)
						})
						continue
					}

					log.Printf("hand shake ok\n")

					m := make(map[resp.Value]resp.Value)
					m[resp.ToBulkString("server")] = resp.ToBulkString("redis")
					m[resp.ToBulkString("version")] = resp.ToBulkString("6.0.16")
					m[resp.ToBulkString("proto")] = resp.ToInteger(3)
					m[resp.ToBulkString("id")] = resp.ToInteger(command.ConnectionID)
					m[resp.ToBulkString("mode")] = resp.ToBulkString("standalone")
					m[resp.ToBulkString("role")] = resp.ToBulkString("master")
					m[resp.ToBulkString("modules")] = resp.ToArray()

					toClientHello := resp.ToMap(m)
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, toClientHello)
					})
				case *commands.CommandCommand:
					m := resp.ToArray()
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, m)
					})
				case *commands.UnknownCommand:
					be, err := resp.ToSimpleError("ERR", "unkonown or wrong command")
					if err != nil {
						panic(err)
					}
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, be)
					})
				case *commands.SetCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					if (innerCmd.Nx && obj != nil) || (innerCmd.Xx && obj == nil) {
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToNull())
						})
						continue
					}

					var ttl *time.Time
					if innerCmd.Px != nil {
						ttl_ := time.Now()
						ttl_ = ttl_.Add(time.Millisecond * time.Duration(*innerCmd.Px))
						ttl = &ttl_
					}
					if innerCmd.Ex != nil {
						ttl_ := time.Now()
						ttl_ = ttl_.Add(time.Second * time.Duration(*innerCmd.Ex))
						ttl = &ttl_
					}

					if obj == nil {
						obj = new(datastructure.RdsObject)
					}

					var data interface{}
					if innerCmd.Encoding == datastructure.EncodingStringRaw {
						data = &datastructure.StringRaw{
							Data: innerCmd.ValueRaw,
						}
					} else {
						data = &datastructure.StringInt{
							Data: innerCmd.ValueInt,
						}
					}
					mem.bigKV[innerCmd.Key] = &datastructure.RdsObject{
						Type:     datastructure.TypeString,
						Encoding: innerCmd.Encoding,
						Data:     data,
						TTL:      ttl,
					}
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToBulkString("OK"))
					})
				case *commands.TTLCommand:
					obj, now := mem.GetRdsObj(innerCmd.Key)
					if obj == nil {
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(-2))
						})
						continue
					}
					if obj.TTL == nil {
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(-1))
						})
						continue
					}
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(int64((*obj.TTL).Sub(now).Seconds())))
					})
				}
			}
		} else {
			mem.cond.Wait()
			mem.mu.Unlock()
		}
	}
}
