package main

import (
	"encoding/json"
	"fmt"
	io "github.com/izqui/go-socket.io"
	"github.com/izqui/helpers"
	"net/http"
)

type SocketCallback chan io.Socket
type DataCallback chan string

type SocketServer struct {
	ConnectCallback                             SocketCallback
	NodeCallback, LinkCallback, MessageCallback DataCallback
	Sockets                                     []*io.Socket

	Server *io.Server
}

func setupWebSocket() *SocketServer {

	socket := &SocketServer{}
	socket.Server = io.NewServer(io.DefaultConfig)

	http.Handle("/socket.io/", socket.Server)
	http.Handle("/", http.FileServer(http.Dir("./visualization/public")))

	return socket
}

func (s *SocketServer) Listen(port string) {

	callbackFunction := func(cb DataCallback) func(d string) {
		return func(d string) {

			if cb != nil {
				cb <- d
			}
		}
	}

	socket.Server.On("connection", func(so io.Socket) {

		s.Sockets = append(s.Sockets, &so)

		//Not working for some reason
		so.On("addlink", callbackFunction(socket.LinkCallback))
		so.On("message", callbackFunction(socket.MessageCallback))

		socket.ConnectCallback <- so
	})

	socket.Server.On("addnode", func(a string) {

		fmt.Println("Add node")
		go BootUpNode(helpers.RandomString(5), 0)
		go BootUpNode(helpers.RandomString(5), 0)
		go BootUpNode(helpers.RandomString(5), 0)
		go BootUpNode(helpers.RandomString(5), 0)
		fmt.Println("Nodes")
		return

	})

	socket.Server.On("error", func(so io.Socket, err error) {
		fmt.Println("error:", err)
	})

	fmt.Println("Interface running on port", port)
	panic(http.ListenAndServe(":"+port, nil))
}

func (s *SocketServer) SendNodes(so io.Socket, nodes ...*Node) {

	var sendnodes func(ttl int, so io.Socket, nodes ...*Node)
	sendnodes = func(ttl int, so io.Socket, nodes ...*Node) {

		//If destination socket is nil, broadcast to all connected sockets
		if so == nil {

			//Avoid entering in loops when sockets have disconnected
			if ttl != 0 {
				return
			}

			for _, so := range s.Sockets {

				sendnodes(1, *so, nodes...)
			}

		} else {

			for _, node := range nodes {

				j, _ := json.Marshal(node)
				so.Emit("addnode", string(j))
			}
		}
	}
	sendnodes(0, so, nodes...)
}
