package handler

import (
	"log"

	goreactor "github.com/markity/go-reactor"
	eventloop "github.com/markity/go-reactor/pkg/event_loop"
)

type connManager map[int64]goreactor.TCPConnection

func DoOnLoop(loop eventloop.EventLoop) {
	log.Printf("loop %v setup", loop.GetID())
	connMap := make(connManager)
	loop.SetContext("conn", connMap)
}
