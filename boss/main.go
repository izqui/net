package main

import (
	"flag"
	"fmt"
)

var (
	port = flag.String("port", "3000", "boss port")
)

var nodes = []*Node{}

func init() {

	flag.Parse()
}
func main() {

	listener := setupTCPListener(*port)
	fmt.Println("TCP connection opened on", *port)

	BootUpNode("A", 0)
	BootUpNode("B", 0)

	cb := make(ConnectionCallback)
	go listenTCP(listener, cb)

	for {

		select {

		case conn := <-cb:

			fmt.Println("New connection")

			node := &Node{BossConnection: conn}
			nodes = append(nodes, node)

			go node.ListenForConnections()
			node.GetInfo()
		}
	}
}

func panicOnError(err error) {

	if err != nil {

		panic(err)
	}
}
