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
			return nil, errors.New("unexpected ping param number")
		}

		second := getExactlyNthBulkStringInArray(val, 1)
		if second == nil {
			return nil, errors.New("unexpected ping second param")
		}

		return &PingCommand{
			Message: second,
		}, nil
	case "client":
		if len(val.Data) != 4 {
			return &UnknownCommand{}, nil
		}

		second := pointerStringToLower(getExactlyNthBulkStringInArray(val, 1))
		if second == nil {
			return nil, errors.New("unexpected")
		}
		if *second != "setinfo" {
			return nil, errors.New("unexpected")
		}

		third := getExactlyNthBulkStringInArray(val, 2)
		if third == nil {
			return nil, errors.New("client thrid param type is unexpected")
		}

		forth := getExactlyNthBulkStringInArray(val, 3)
		if forth == nil {
			return nil, errors.New("client forth param type is unexpected")
		}

		switch *third {
		case "LIB-NAME":
			return &SetInfoLibNameCommand{
				LibName: *forth,
			}, nil
		case "LIB-VER":
			return &SetInfoLibVersionCommand{
				LibVersion: *forth,
			}, nil
		default:
			return nil, errors.New("client thrid param is unexpected")
		}
	case "command":
		if len(val.Data) != 1 {
			return nil, errors.New("unexpected command number param number")
		}

		return &CommandCommand{}, nil
	case "set":
		// set k v [EX seconds|PX milliseconds] [NX|XX]
		if len(val.Data) < 3 {
			return &UnknownCommand{}, nil
		}

		second_ := getExactlyNthBulkStringInArray(val, 1)
		if second_ == nil {
			return nil, errors.New("set command: unexpected second param type")
		}

		third_ := getExactlyNthBulkStringInArray(val, 2)
		if third_ == nil {
			return nil, errors.New("set command: unexpected thrid param type")
		}

		// key
		second := *second_
		// value
		third := *third_

		isInt64 := true

		// check if value can be contain by int64

		// start with 0
		if len(third) > 1 && third[0] == '0' {
			isInt64 = false
		}

		// has + char
		if len(third) > 0 && third[0] == '+' {
			isInt64 = false
		}

		// check is int64
		i64, err := strconv.ParseInt(third, 10, 64)
		if err != nil {
			isInt64 = false
		}

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
				entry_ := pointerStringToLower(getExactlyNthBulkStringInArray(val, i))
				if entry_ == nil {
					return nil, errors.New("unexpected")
				}
				entry := *entry_
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
					exNumberStr_ := getExactlyNthBulkStringInArray(val, i+1)
					if exNumberStr_ == nil {
						return nil, errors.New("unexpected")
					}

					exNumberStr := *exNumberStr_
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
					pxNumberStr_ := getExactlyNthBulkStringInArray(val, i+1)
					if pxNumberStr_ == nil {
						return nil, errors.New("unexpected")
					}

					pxNumberStr := *pxNumberStr_
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
	default:
		return &UnknownCommand{}, nil
	}
}

func NewCommandsParser() *CommandParser {
	return &CommandParser{}
}
