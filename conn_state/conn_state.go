package connstate

import (
	"go-rds/pkg/queue"
	"go-rds/pkg/resp"
)

type ConnState struct {
	HandleShake bool
	RespParser  resp.RespParser
	ValueToken  queue.Queue
	ConnID      int64
}
