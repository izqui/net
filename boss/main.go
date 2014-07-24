package main

import (
	"flag"
	"fmt"
	"github.com/izqui/helpers"
	_ "time"
)

var (
	port          = flag.String("port", "3000", "boss port")
	interfacePort = flag.String("interface", "7777", "interface port")
)

var nodes = []*Node{}
var socket *SocketServer

func init() {

	flag.Parse()
}
func main() {

	socket = setupWebSocket()
	socket.ConnectCallback = make(SocketCallback)
	socket.NodeCallback = make(DataCallback, 10000)
	socket.MessageCallback = make(DataCallback)
	socket.LinkCallback = make(DataCallback)

	socket.OnNode = func(a string) {

		go BootUpNode(helpers.RandomString(5), 0)
		fmt.Println("finish")
	}

	go socket.Listen(*interfacePort)

	listener := setupTCPListener(*port)
	fmt.Println("TCP connection opened on", *port)
	bossCb := make(ConnectionCallback)
	go listenTCP(listener, bossCb)

	cb2 := make(chan bool)
	go input(cb2)

	for {

		select {

		case conn := <-bossCb:

			fmt.Println("New peer connection")

			node := &Node{BossConnection: conn}
			nodes = append(nodes, node)

			go node.ListenForConnections(func() {

				fmt.Println("Node disconnected")
			})
			node.GetInfo()

		case so := <-socket.ConnectCallback:

			fmt.Println("New socket connection")
			go socket.SendNodes(so, nodes...)

		case <-socket.NodeCallback:
			//Not working for some reason
			go BootUpNode(helpers.RandomString(5), 0)

		case <-cb2:

			//On enter. Just a test
			socket.NodeCallback <- "hello"

			/*fmt.Println(nodes[0].Id)
			nodes[1].ConnectToNode(nodes[0].PeerAddress)
			time.Sleep(1 * time.Second)
			nodes[1].SendMessageToNode(nodes[0].Id)*/
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
