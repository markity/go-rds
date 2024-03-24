package handler

import (
	connstate "go-rds/conn_state"
	"go-rds/global"
	"go-rds/pkg/queue"
	"go-rds/pkg/resp"
	"go-rds/types"

	goreactor "github.com/markity/go-reactor"
)

func OnConnection(t goreactor.TCPConnection) {
	if t.IsConnected() {
		state := &connstate.ConnState{
			HandleShake: false,
			RespParser:  resp.NewRespParser(),
			ValueToken:  make(queue.Queue, 0),
			ConnID:      global.ConnIDCounter.Add(1),
		}
		t.SetContext("state", state)
		t.GetEventLoop().MustGetContext("conn").(types.ConnManager)[state.ConnID] = t
	} else {
		state := t.MustGetContext("state").(*connstate.ConnState)
		loop := t.GetEventLoop()
		conns := loop.MustGetContext("conn").(types.ConnManager)
		delete(conns, state.ConnID)
	}
}
