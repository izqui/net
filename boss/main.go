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

	go BootUpNode("A", 0)
	go BootUpNode("B", 0)

	cb := make(ConnectionCallback)
	go listenTCP(listener, cb)

	cb2 := make(chan bool)
	go input(cb2)

	for {

		select {

		case conn := <-cb:

			fmt.Println("New connection")

			node := &Node{BossConnection: conn}
			nodes = append(nodes, node)

			go node.ListenForConnections()
			node.GetInfo()

		case <-cb2:

			nodes[1].ConnectToNode(nodes[0].PeerAddress)
		}
	}
}

func input(a chan bool) {

	for {

		fmt.Scanf("\n")
		a <- true
	}
}
func panicOnError(err error) {

	if err != nil {

		panic(err)
	}
}
