package memoryloop

import (
	"go-rds/pkg/resp"
	"go-rds/types"

	eventloop "github.com/markity/go-reactor/pkg/event_loop"
)

func CallbackFunc(loop eventloop.EventLoop, connID int64, response resp.Value) {
	conn, ok := loop.MustGetContext("conn").(types.ConnManager)[connID]
	if !ok {
		return
	}

	bs, err := response.EncodeToBytes()
	if err != nil {
		panic(err)
	}
	conn.Send(bs)
}
