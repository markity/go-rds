package types

import (
	"go-rds/pkg/resp"
)

type ConnState struct {
	RespParser resp.RespParser
	ConnID     int64
}
