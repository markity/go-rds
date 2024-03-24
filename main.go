package main

import (
	"go-rds/handler"

	goreactor "github.com/markity/go-reactor"
	eventloop "github.com/markity/go-reactor/pkg/event_loop"
)

func main() {
	// new loop
	loop := eventloop.NewEventLoop()

	// new tcp server
	server := goreactor.NewTCPServer(loop, "127.0.0.1:6379", 16, goreactor.RoundRobin())

	// get all io loops setup
	_, ioLoops := server.GetAllLoops()
	for _, ioLoop := range ioLoops {
		ioLoop.DoOnLoop(handler.DoOnLoop)
	}

	// set handlers
	server.SetConnectionCallback(handler.OnConnection)
	server.SetMessageCallback(handler.OnMessage)

	// run server, register related epoll fds
	err := server.Start()
	if err != nil {
		panic(err)
	}

	// runs event loop
	loop.Loop()
}
