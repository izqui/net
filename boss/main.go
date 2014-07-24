package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/izqui/helpers"
	_ "time"
)

var (
	port          = flag.String("port", "3000", "boss port")
	interfacePort = flag.String("interface", "7777", "interface port")
)

var nodes = NodeSlice{}
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

			go BootUpNode(helpers.RandomString(5), 0)

		case l := <-socket.LinkCallback:

			link := new(Link)
			json.Unmarshal([]byte(l), link)

			s := nodes.FindNode(link.Source)
			d := nodes.FindNode(link.Destination)

			go s.ConnectToNode(d.PeerAddress)

		case m := <-socket.MessageCallback:

			message := new(BossMessage)
			json.Unmarshal([]byte(m), message)

			s := nodes.FindNode(message.From)
			d := nodes.FindNode(message.To)

			go s.SendMessageToNode(d.Id)

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
