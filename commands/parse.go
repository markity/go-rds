package commands

import (
	"errors"
	datastructure "go-rds/data_structure"
	"go-rds/pkg/resp"
	"strconv"
	"strings"
)

type CommandParser struct {
}

func getExactlyNthBulkStringInArray(val *resp.Array, n int) *string {
	if n >= len(val.Data) {
		return nil
	}

	bs_, ok := val.Data[n].(*resp.BulkString)
	if !ok {
		return nil
	}

	return &bs_.Data
}

func pointerStringToLower(str *string) *string {
	if str == nil {
		return nil
	}

	*str = strings.ToLower(*str)
	return str
}

func (*CommandParser) Parse(val_ resp.Value) (interface{}, error) {
	val, ok := val_.(*resp.Array)
	if !ok {
		return nil, errors.New("parse error, unexpected value type")
	}
	if len(val.Data) == 0 {
		return nil, errors.New("parse error, array length is 0")
	}

	for _, v := range val.Data {
		_, ok := v.(*resp.BulkString)
		if !ok {
			return nil, errors.New("command params must be bulk string")
		}
	}

	first := pointerStringToLower(getExactlyNthBulkStringInArray(val, 0))
	if first == nil {
		return nil, errors.New("need first bulk string as base command")
	}
	switch *first {
	case "hello":
		second := getExactlyNthBulkStringInArray(val, 1)
		if second == nil {
			return &UnknownCommand{}, nil
		}

		return &HelloCommand{
			Proto: *second,
		}, nil
	case "ping":
		if len(val.Data) == 1 {
			return &PingCommand{
				Message: nil,
			}, nil
		}

		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		second := getExactlyNthBulkStringInArray(val, 1)
		if second == nil {
			panic("unexpected")
		}

		return &PingCommand{
			Message: second,
		}, nil
	case "echo":
		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		// no panic here
		second := *getExactlyNthBulkStringInArray(val, 1)

		return &EchoCommand{
			Message: second,
		}, nil
	case "client":
		if len(val.Data) != 4 {
			return &UnknownCommand{}, nil
		}

		second := *pointerStringToLower(getExactlyNthBulkStringInArray(val, 1))
		if second != "setinfo" {
			return &UnknownCommand{}, nil
		}

		third := *getExactlyNthBulkStringInArray(val, 2)

		forth := *getExactlyNthBulkStringInArray(val, 3)

		switch third {
		case "LIB-NAME":
			return &SetInfoLibNameCommand{
				LibName: forth,
			}, nil
		case "LIB-VER":
			return &SetInfoLibVersionCommand{
				LibVersion: forth,
			}, nil
		default:
			return &UnknownCommand{}, nil
		}
	case "command":
		if len(val.Data) != 1 {
			return &UnknownCommand{}, nil
		}

		return &CommandCommand{}, nil
	case "set":
		// set k v [EX seconds|PX milliseconds] [NX|XX]
		if len(val.Data) < 3 {
			return &UnknownCommand{}, nil
		}

		// key
		second := *getExactlyNthBulkStringInArray(val, 1)
		// value
		third := *getExactlyNthBulkStringInArray(val, 2)

		i64, isInt64 := ToolIsStringInteger(third)

		encoding := datastructure.EncodingStringRaw
		if isInt64 {
			encoding = datastructure.EncodingStringInt
		}

		xx := false
		nx := false
		var ex *int
		var px *int
		if len(val.Data) > 3 {
			for i := 3; i < len(val.Data); i++ {
				// no panic here
				entry := *pointerStringToLower(getExactlyNthBulkStringInArray(val, i))
				switch entry {
				case "xx":
					if nx {
						return &UnknownCommand{}, nil
					}

					xx = true
				case "nx":
					if xx {
						return &UnknownCommand{}, nil
					}
					nx = true
				case "ex":
					if px != nil {
						return &UnknownCommand{}, nil
					}

					if i+1 >= len(val.Data) {
						return &UnknownCommand{}, nil
					}
					exNumberStr := *getExactlyNthBulkStringInArray(val, i+1)
					if len(exNumberStr) == 0 {
						return nil, errors.New("unexpected")
					}

					if exNumberStr[0] == '+' || exNumberStr[0] == '-' {
						return &UnknownCommand{}, nil
					}

					i64, err := strconv.ParseInt(exNumberStr, 10, 32)
					if err != nil {
						return &UnknownCommand{}, nil
					}
					i32 := int(i64)
					ex = &i32
					i++
				case "px":
					if ex != nil {
						return &UnknownCommand{}, nil
					}

					if i+1 >= len(val.Data) {
						return &UnknownCommand{}, nil
					}
					pxNumberStr := *getExactlyNthBulkStringInArray(val, i+1)
					if len(pxNumberStr) == 0 {
						return nil, errors.New("unexpected")
					}

					if pxNumberStr[0] == '+' || pxNumberStr[0] == '-' {
						return &UnknownCommand{}, nil
					}

					i64, err := strconv.ParseInt(pxNumberStr, 10, 32)
					if err != nil {
						return &UnknownCommand{}, nil
					}
					i32 := int(i64)
					px = &i32
					i++
				default:
					return &UnknownCommand{}, nil
				}
			}
		}

		return &SetCommand{
			Key:      second,
			Encoding: encoding,
			ValueRaw: third,
			ValueInt: i64,
			Nx:       nx,
			Xx:       xx,
			Ex:       ex,
			Px:       px,
		}, nil
	case "ttl":
		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		// no panic here
		key := *getExactlyNthBulkStringInArray(val, 1)
		return &TTLCommand{
			Key: key,
		}, nil
	case "pttl":
		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		// no panic here
		key := *getExactlyNthBulkStringInArray(val, 1)
		return &PTTLCommand{
			Key: key,
		}, nil
	case "get":
		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		// no panic here
		key := *getExactlyNthBulkStringInArray(val, 1)
		return &GetCommand{
			Key: key,
		}, nil
	case "append":
		if len(val.Data) != 3 {
			return &UnknownCommand{}, nil
		}

		key := *getExactlyNthBulkStringInArray(val, 1)
		value := *getExactlyNthBulkStringInArray(val, 2)
		return &AppendCommand{
			Key:     key,
			Content: value,
		}, nil
	case "incr":
		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		key := *getExactlyNthBulkStringInArray(val, 1)
		return &IncrCommand{
			Key: key,
		}, nil
	case "incrby":
		if len(val.Data) != 3 {
			return &UnknownCommand{}, nil
		}

		key := *getExactlyNthBulkStringInArray(val, 1)
		incr, ok := ToolIsStringInteger(*getExactlyNthBulkStringInArray(val, 2))
		if !ok {
			return &UnknownCommand{}, nil
		}

		return &IncrByCommand{
			Key: key,
			By:  incr,
		}, nil
	case "decr":
		if len(val.Data) != 2 {
			return &UnknownCommand{}, nil
		}

		key := *getExactlyNthBulkStringInArray(val, 1)
		return &DecrCommand{
			Key: key,
		}, nil
	case "decrby":
		if len(val.Data) != 3 {
			return &UnknownCommand{}, nil
		}

		key := *getExactlyNthBulkStringInArray(val, 1)
		decr, ok := ToolIsStringInteger(*getExactlyNthBulkStringInArray(val, 2))
		if !ok {
			return &UnknownCommand{}, nil
		}

		return &DecrByCommand{
			Key: key,
			By:  decr,
		}, nil
	case "del":
		if len(val.Data) == 1 {
			return &UnknownCommand{}, nil
		}

		var keys []string
		for i := 1; i < len(val.Data); i++ {
			k := *getExactlyNthBulkStringInArray(val, i)
			keys = append(keys, k)
		}
		return &DelCommand{
			Key: keys,
		}, nil
	default:
		return &UnknownCommand{}, nil
	}
}

func NewCommandsParser() *CommandParser {
	return &CommandParser{}
}
