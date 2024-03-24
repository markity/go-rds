package handler

import (
	"go-rds/types"
	"log"

	eventloop "github.com/markity/go-reactor/pkg/event_loop"
)

func DoOnLoop(loop eventloop.EventLoop) {
	log.Printf("loop %v setup", loop.GetID())
	connMap := make(types.ConnManager)
	loop.SetContext("conn", connMap)
}
