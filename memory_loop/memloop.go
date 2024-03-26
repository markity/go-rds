package memoryloop

import (
	"fmt"
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
				case *commands.EchoCommand:
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToBulkString(innerCmd.Message))
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
				case *commands.PTTLCommand:
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
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(int64((*obj.TTL).Sub(now).Milliseconds())))
					})
				case *commands.GetCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					if obj == nil {
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToNull())
						})
						continue
					}

					if obj.Type != datastructure.TypeString {
						se, err := resp.ToSimpleError("WRONGTYPE", "Operation against a key holding the wrong kind of value")
						if err != nil {
							panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, se)
						})
						continue
					}

					var data string
					switch obj.Encoding {
					case datastructure.EncodingStringRaw:
						data = obj.Data.(*datastructure.StringRaw).Data
					case datastructure.EncodingStringInt:
						data = fmt.Sprint(obj.Data.(*datastructure.StringRaw).Data)
					}

					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToBulkString(data))
					})
				case *commands.AppendCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					// 如果不存在, 那么直接创建
					if obj == nil {
						obj := new(datastructure.RdsObject)
						obj.Type = datastructure.TypeString
						i64, isInt64 := commands.ToolIsStringInteger(innerCmd.Content)
						if isInt64 {
							obj.Encoding = datastructure.EncodingStringInt
							obj.Data = &datastructure.StringInt{
								Data: i64,
							}
						} else {
							obj.Encoding = datastructure.EncodingStringRaw
							obj.Data = &datastructure.StringRaw{
								Data: innerCmd.Content,
							}
						}
						obj.TTL = nil
						mem.SetObj(innerCmd.Key, obj)
						length := len(innerCmd.Content)
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(int64(length)))
						})
						continue
					}

					// 如果key存在, 但是类型不对
					if obj.Type != datastructure.TypeString {
						be, err := resp.ToSimpleError("ERR", "unkonown or wrong command")
						if err != nil {
							panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, be)
						})
						continue
					}

					// 如果类型对
					var originStr string
					if obj.Encoding == datastructure.EncodingStringInt {
						originStr = fmt.Sprint(obj.Data.(*datastructure.StringInt).Data)
					} else {
						originStr = obj.Data.(*datastructure.StringRaw).Data
					}
					newStr := originStr + innerCmd.Content

					i64, isInt64 := commands.ToolIsStringInteger(newStr)
					if isInt64 {
						obj.Encoding = datastructure.EncodingStringInt
						obj.Data = &datastructure.StringInt{
							Data: i64,
						}
					} else {
						obj.Encoding = datastructure.EncodingStringRaw
						obj.Data = &datastructure.StringRaw{
							Data: newStr,
						}
					}
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(int64(len(newStr))))
					})
				case *commands.IncrCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					// 如果不存在, 那么直接创建
					if obj == nil {
						obj := new(datastructure.RdsObject)
						obj.Type = datastructure.TypeString
						obj.Encoding = datastructure.EncodingStringInt
						obj.Data = &datastructure.StringInt{
							Data: 1,
						}
						obj.TTL = nil
						mem.SetObj(innerCmd.Key, obj)
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(1))
						})
						continue
					}

					if obj.Encoding != datastructure.EncodingStringInt {
						se, err := resp.ToSimpleError("ERR", "value is not an integer")
						if err != nil {
							log.Panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, se)
						})
						continue
					}

					si := obj.Data.(*datastructure.StringInt)
					si.Data++
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(si.Data))
					})
				case *commands.IncrByCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					// 如果不存在, 那么直接创建
					if obj == nil {
						obj := new(datastructure.RdsObject)
						obj.Type = datastructure.TypeString
						obj.Encoding = datastructure.EncodingStringInt
						obj.Data = &datastructure.StringInt{
							Data: innerCmd.By,
						}
						obj.TTL = nil
						mem.SetObj(innerCmd.Key, obj)
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(innerCmd.By))
						})
						continue
					}

					if obj.Encoding != datastructure.EncodingStringInt {
						se, err := resp.ToSimpleError("ERR", "value is not an integer")
						if err != nil {
							log.Panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, se)
						})
						continue
					}

					si := obj.Data.(*datastructure.StringInt)
					si.Data += innerCmd.By
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(si.Data))
					})
				case *commands.DecrCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					// 如果不存在, 那么直接创建
					if obj == nil {
						obj := new(datastructure.RdsObject)
						obj.Type = datastructure.TypeString
						obj.Encoding = datastructure.EncodingStringInt
						obj.Data = &datastructure.StringInt{
							Data: -1,
						}
						obj.TTL = nil
						mem.SetObj(innerCmd.Key, obj)
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(-1))
						})
						continue
					}

					if obj.Encoding != datastructure.EncodingStringInt {
						se, err := resp.ToSimpleError("ERR", "value is not an integer")
						if err != nil {
							log.Panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, se)
						})
						continue
					}

					si := obj.Data.(*datastructure.StringInt)
					si.Data--
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(si.Data))
					})
				case *commands.DecrByCommand:
					obj, _ := mem.GetRdsObj(innerCmd.Key)
					// 如果不存在, 那么直接创建
					if obj == nil {
						obj := new(datastructure.RdsObject)
						obj.Type = datastructure.TypeString
						obj.Encoding = datastructure.EncodingStringInt
						obj.Data = &datastructure.StringInt{
							Data: -innerCmd.By,
						}
						obj.TTL = nil
						mem.SetObj(innerCmd.Key, obj)
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(-innerCmd.By))
						})
						continue
					}

					if obj.Encoding != datastructure.EncodingStringInt {
						se, err := resp.ToSimpleError("ERR", "value is not an integer")
						if err != nil {
							log.Panic(err)
						}
						command.Loop.RunInLoop(func() {
							CallbackFunc(command.Loop, command.ConnectionID, se)
						})
						continue
					}

					si := obj.Data.(*datastructure.StringInt)
					si.Data -= innerCmd.By
					command.Loop.RunInLoop(func() {
						CallbackFunc(command.Loop, command.ConnectionID, resp.ToInteger(si.Data))
					})
				}
			}
		} else {
			mem.cond.Wait()
			mem.mu.Unlock()
		}
	}
}
