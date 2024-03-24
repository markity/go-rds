package memoryloop

import (
	"go-rds/commands"
	"go-rds/pkg/resp"
	"go-rds/types"
)

func MemLoopCallback(cmd *commands.MemLoopCommandCover) {
	conn, ok := cmd.Loop.MustGetContext("conn").(types.ConnManager)[cmd.ConnectionID]
	if !ok {
		return
	}

	switch cmd.Command.(type) {
	case *commands.PingCommand:
		bs, err := resp.ToBulkString("PONG").EncodeToBytes()
		if err != nil {
			panic(err)
		}
		conn.Send(bs)
	case *commands.SetInfoLibNameCommand:
		ss, err := resp.ToSimpleString("OK")
		if err != nil {
			panic(err)
		}
		bs, err := ss.EncodeToBytes()
		if err != nil {
			panic(err)
		}
		conn.Send(bs)
	case *commands.SetInfoLibVersionCommand:
		ss, err := resp.ToSimpleString("OK")
		if err != nil {
			panic(err)
		}
		bs, err := ss.EncodeToBytes()
		if err != nil {
			panic(err)
		}
		conn.Send(bs)
	}
}
