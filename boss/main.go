package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	port          = flag.String("port", "3000", "boss port")
	interfacePort = flag.String("interface", "7777", "interface port")
)

var nodes = []*Node{}
var socket Socket

func init() {

	flag.Parse()
}
func main() {

	socket = setupWebSocket()
	go socket.Listen(*interfacePort)

	//go BootUpNode("A", 0)
	//go BootUpNode("B", 0)

	listener := setupTCPListener(*port)
	fmt.Println("TCP connection opened on", *port)
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

			fmt.Println(nodes[0].Id)
			nodes[1].ConnectToNode(nodes[0].PeerAddress)
			time.Sleep(1 * time.Second)
			nodes[1].SendMessageToNode(nodes[0].Id)
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
