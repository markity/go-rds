package handler

import (
	"fmt"
	connstate "go-rds/conn_state"
	"go-rds/pkg/resp"
	"log"

	goreactor "github.com/markity/go-reactor"
	"github.com/markity/go-reactor/pkg/buffer"
)

func OnMessage(conn goreactor.TCPConnection, buffer buffer.Buffer) {
	state_, _ := conn.GetContext("state")
	state := state_.(*connstate.ConnState)
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
		state.ValueToken.PushBack(value)
	}

	if !state.HandleShake {
		if state.ValueToken.Size() == 0 {
			return
		}

		handShakeValue := state.ValueToken.PopFront()
		handShakeList, ok := handShakeValue.(*resp.Array)
		if !ok {
			log.Printf("hand shake failed\n")
			conn.ForceClose()
			return
		}

		if len(handShakeList.Data) != 2 {
			log.Printf("hand shake failed\n")
			conn.ForceClose()
			return
		}

		helloMsg, ok1 := handShakeList.Data[0].(*resp.BulkString)
		versionMsg, ok2 := handShakeList.Data[1].(*resp.BulkString)
		if !ok1 || !ok2 || helloMsg.Data != "hello" || versionMsg.Data != "3" {
			log.Printf("hand shake failed\n")
			conn.ForceClose()
			return
		}

		log.Printf("hand shake ok\n")

		m := make(map[resp.Value]resp.Value)
		m[resp.ToBulkString("server")] = resp.ToBulkString("redis")
		m[resp.ToBulkString("version")] = resp.ToBulkString("6.0.16")
		m[resp.ToBulkString("proto")] = resp.ToInteger(3)
		m[resp.ToBulkString("id")] = resp.ToInteger(state.ConnID)
		m[resp.ToBulkString("mode")] = resp.ToBulkString("standalone")
		m[resp.ToBulkString("role")] = resp.ToBulkString("master")
		m[resp.ToBulkString("modules")] = resp.ToArray()
		toClientHello := resp.ToMap(m)
		toClientHelloBytes, err := toClientHello.EncodeToBytes()
		if err != nil {
			panic(err)
		}
		conn.Send(toClientHelloBytes)
		fmt.Printf("%s\n", toClientHelloBytes)
		state.HandleShake = true
	}

	// 消费数据
}
