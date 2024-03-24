package handler

import (
	connstate "go-rds/conn_state"
	"go-rds/global"
	"go-rds/pkg/queue"
	"go-rds/pkg/resp"

	goreactor "github.com/markity/go-reactor"
)

func OnConnection(t goreactor.TCPConnection) {
	if t.IsConnected() {
		t.SetContext("state", &connstate.ConnState{
			HandleShake: false,
			RespParser:  resp.NewRespParser(),
			ValueToken:  make(queue.Queue, 0),
			ConnID:      global.ConnIDCounter.Add(1),
		})
	} else {
		state := t.MustGetContext("state").(*connstate.ConnState)
		loop := t.GetEventLoop()
		conns := loop.MustGetContext("conn").(connManager)
		delete(conns, state.ConnID)
	}
}
