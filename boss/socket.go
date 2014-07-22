package main

import (
	"fmt"
	io "github.com/googollee/go-socket.io"
	"net/http"
)

type Socket struct {
	server *io.SocketIOServer
}

func setupWebSocket() Socket {

	conf := &io.Config{ClosingTimeout: 2}

	s := io.NewSocketIOServer(conf)
	s.Handle("/", http.FileServer(http.Dir("./public")))

	return Socket{server: s}
}

func (s Socket) Listen(port string) {

	fmt.Println("Interface running on port", port)
	panic(http.ListenAndServe(":"+port, s.server))
}

func (s Socket) Broadcast(name string, args ...interface{}) {

	s.server.Broadcast(name, args...)
}
