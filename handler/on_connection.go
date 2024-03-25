package handler

import (
	"go-rds/global"
	"go-rds/pkg/resp"
	"go-rds/types"
	"log"

	goreactor "github.com/markity/go-reactor"
)

func OnConnection(t goreactor.TCPConnection) {
	t.SetDisConnectedCallback(func(t goreactor.TCPConnection) {
		log.Printf("disconnected\n")
		state := t.MustGetContext("state").(*types.ConnState)
		loop := t.GetEventLoop()
		conns := loop.MustGetContext("conn").(types.ConnManager)
		delete(conns, state.ConnID)
	})

	log.Printf("connected\n")
	state := &types.ConnState{
		RespParser: resp.NewRespParser(),
		ConnID:     global.ConnIDCounter.Add(1),
	}
	t.SetContext("state", state)
	t.GetEventLoop().MustGetContext("conn").(types.ConnManager)[state.ConnID] = t
}
