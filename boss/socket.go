package main

import (
	"encoding/json"
	"fmt"
	io "github.com/izqui/go-socket.io"
	_ "github.com/izqui/helpers"
	"net/http"
)

type SocketCallback chan io.Socket
type DataCallback chan string

type Link struct {
	Source      string `json:"source"`
	Destination string `json:"target"`
}

type SocketServer struct {
	ConnectCallback                             SocketCallback
	NodeCallback, LinkCallback, MessageCallback DataCallback
	OnNode                                      func(a string)
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
		so.On("disconnection", func() {

			fmt.Println("disconnect")
			so = nil
		})

		so.On("addnode", callbackFunction(socket.NodeCallback))
		so.On("addlink", callbackFunction(socket.LinkCallback))
		so.On("message", callbackFunction(socket.MessageCallback))

		socket.ConnectCallback <- so
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

			for _, node := range nodes {

				s.SendLinks(so, node.GetLinks()...)
			}
		}
	}
	sendnodes(0, so, nodes...)
}

func (s *SocketServer) SendLinks(so io.Socket, links ...Link) {

	var sendlinks func(ttl int, so io.Socket, links ...Link)
	sendlinks = func(ttl int, so io.Socket, links ...Link) {

		//If destination socket is nil, broadcast to all connected sockets
		if so == nil {

			//Avoid entering in loops when sockets have disconnected
			if ttl != 0 {
				return
			}

			for _, so := range s.Sockets {

				sendlinks(1, *so, links...)
			}

		} else {

			for _, link := range links {

				j, _ := json.Marshal(link)
				so.Emit("addlink", string(j))
			}
		}
	}
	sendlinks(0, so, links...)
}
