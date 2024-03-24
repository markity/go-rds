package commands

import (
	"errors"
	"go-rds/pkg/resp"
)

type CommandParser struct {
}

func (*CommandParser) Parse(val_ resp.Value) (interface{}, error) {
	val, ok := val_.(*resp.Array)
	if !ok {
		return nil, errors.New("parse error, unexpected value type")
	}

	if len(val.Data) == 0 {
		return nil, errors.New("parse error, array length is 0")
	}

	baseCommand, ok := val.Data[0].(*resp.BulkString)
	if !ok {
		return nil, errors.New("parse base command error, unexpected type")
	}
	first := baseCommand.Data
	switch first {
	case "ping":
		if len(val.Data) == 1 {
			return &PingCommand{
				Message: nil,
			}, nil
		}

		if len(val.Data) != 2 {
			return nil, errors.New("unexpected ping param number")
		}

		str_, ok := val.Data[1].(*resp.BulkString)
		if !ok {
			return nil, errors.New("unexpected ping second param type")
		}

		str := str_.Data

		return &PingCommand{
			Message: &str,
		}, nil
	case "client":
		if len(val.Data) != 4 {
			return nil, errors.New("unexpected client command number param number")
		}

		second_, ok := val.Data[1].(*resp.BulkString)
		if !ok {
			return nil, errors.New("client command second param type is unexpected")
		}
		second := second_.Data
		if second != "setinfo" {
			return nil, errors.New("client second command is not setinfo")
		}

		third_, ok := val.Data[2].(*resp.BulkString)
		if !ok {
			return nil, errors.New("client thrid param type is unexpected")
		}
		third := third_.Data

		forth_, ok := val.Data[3].(*resp.BulkString)
		if !ok {
			return nil, errors.New("client forth param type is unexpected")
		}
		forth := forth_.Data

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
			return nil, errors.New("client thrid param is unexpected")
		}
	default:
		return nil, errors.New("unexpected command")
	}
}

func NewCommandsParser() *CommandParser {
	return &CommandParser{}
}
