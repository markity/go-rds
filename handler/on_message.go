package handler

import (
	"go-rds/commands"
	connstate "go-rds/conn_state"
	"go-rds/global"
	"log"

	goreactor "github.com/markity/go-reactor"
	"github.com/markity/go-reactor/pkg/buffer"
)

func OnMessage(conn goreactor.TCPConnection, buffer buffer.Buffer) {
	state := conn.MustGetContext("state").(*connstate.ConnState)
	err := state.RespParser.Write([]byte(buffer.RetrieveAsString()))
	if err != nil {
		log.Printf("client wrong input: %v", err)
		conn.ForceClose()
		return
	}

	for {
		value, ok := state.RespParser.TakeValue()
		if !ok {
			break
		}

		cmd, err := global.CommandParser.Parse(value)
		if err != nil {
			log.Printf("command parser error: %v", err)
			conn.ForceClose()
			return
		}

		global.MemLoop.QueueCommand(&commands.MemLoopCommandCover{
			Loop:         conn.GetEventLoop(),
			ConnectionID: state.ConnID,
			Command:      cmd,
		})
	}
}
