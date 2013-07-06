package server

import (
	"fmt"
	"github.com/Craig-Macomber/election/msg"
	"net"
)

type HandlerMap map[msg.Type]func(net.Conn)

func Start(handlers HandlerMap) {
	listener, err := net.Listen("tcp", msg.Service)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accept:", err.Error())
			return
		}
		go handelCon(conn, handlers)
	}
}

func ConnectionError(c net.Conn) {
	fmt.Println("connectionError")
	SendBlock(msg.InvalidRequest, []byte("Error!"), c)
}

func handelCon(c net.Conn, handlers HandlerMap) {
	defer c.Close()

	messageType, err := msg.ReadType(c)
	if err != nil {
		fmt.Println("server error reading messageType:", err)
		ConnectionError(c)
	}
	handler, ok := handlers[messageType]
	if !ok {
		fmt.Println("server error missing handler:", messageType)
		ConnectionError(c)
		return
	}

	handler(c)
}

func BlockHandler(f func([]byte, net.Conn), lengthMax uint16) func(net.Conn) {
	return func(c net.Conn) {
		buff, err := msg.ReadBlock(c, lengthMax)
		if err != nil {
			fmt.Println("failed to send messageType:", err)
			ConnectionError(c)
			return
		}
		f(buff, c)
	}
}

func SendBlock(messageType msg.Type, data []byte, c net.Conn) {
	err := msg.WriteBlock(c, messageType, data)
	if err != nil {
		fmt.Println("failed to send block:", err)
		ConnectionError(c)
		return
	}
}
