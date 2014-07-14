package main

import (
	"flag"
	"fmt"
	"github.com/izqui/helpers"
	"io"
	"net"
	"os"
	"time"
)

var (
	scanAddr = "10.0.5.33"
)

var (
	scanPort = flag.String("scan", "3003", "default port for scanning")
	port     = flag.String("port", "0", "your local port")
	id       = flag.String("id", helpers.SHA1([]byte(helpers.RandomString(5))), "id of the peer for the network")
)

var self *Peer
var sentMessages []string

var currentMessage *Message
var messageState = 0

const (
	START_STATE = iota
	MESSAGE_STATE
	CONNECTION_STATE
)

func init() {

	flag.Parse()

	self = new(Peer)
	self.Id = *id
	self.Address = fmt.Sprintf("%s:%s", myIp(), *port)
}

func main() {

	incomingConnection := setupIncomingConnection(self.Address)

	fmt.Println(self.Id, "listening on", self.Address, "scanning on", *scanPort)
	fmt.Println("Network: ", self)

	inputCb := make(chan []byte)
	connectionCb := make(chan []byte)
	searcherCb := make(chan *net.UDPConn)

	go runReadInput(inputCb)
	go runConnectionInput(incomingConnection, connectionCb)
	go searchPeersOnPort(*scanPort, searcherCb)

	for {
		select {

		case input := <-inputCb:
			go inputHandler(input)

		case input := <-connectionCb:
			go self.HandleIncomingConnection(input)

		case connection := <-searcherCb:
			go self.HandleConnectionFound(connection)
		}
	}
}

func runReadInput(cb chan []byte) {

	for {

		in := readInput(os.Stdin)
		cb <- in[:len(in)-1]
	}
}
func runConnectionInput(connection *net.UDPConn, cb chan []byte) {

	for {

		var buffer []byte = make([]byte, 4096)
		n, addr, err := connection.ReadFromUDP(buffer[0:])
		panicOnError(err)

		if addr != nil {

			cb <- buffer[:n]
		}
	}
}

func searchPeersOnPort(port string, cb chan *net.UDPConn) {

	for {

		//network := []string{scanAddr}
		network := []string{scanAddr, "255.255.255.255"}
		for _, address := range network {

			var add = address + ":" + port
			con, err := pingAddress(add)
			if err == nil && con != nil {

				cb <- con
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func inputHandler(input []byte) {

	switch messageState {
	case 0:

		str := string(input)

		if str == "connect" {

			messageState = CONNECTION_STATE

		} else {

			var dest_id = string(input)
			var next_peer = self.FindNearestPeerToId(dest_id)

			if next_peer != nil {

				currentMessage = &Message{Destination: next_peer.Address, FinalDestinationId: dest_id}
				messageState = MESSAGE_STATE

				fmt.Println("Sending message to", dest_id, "through", next_peer)
			} else {

				fmt.Println("Couldn't find peer with that id")
			}
		}

	case 1:

		messageState = START_STATE
		currentMessage.Body = string(input)
		currentMessage.Origin = self
		currentMessage.AssignRandomID()

		self.SendMessage(currentMessage, currentMessage.Destination)

	case 2:

		messageState = START_STATE
		self.StablishConnection(string(input))
	}
}

func panicOnError(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}
