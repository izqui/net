package main

import (
	"flag"
	"fmt"
	"github.com/izqui/helpers"
	"io"
	"net"
	"time"
)

var (
	myIp, scanAddr = "::1", "::1"
)

var (
	scanPort = flag.String("scan", "3003", "default port for scanning")
	port     = flag.String("port", "9999", "your local port")
	id       = flag.String("id", helpers.SHA1([]byte(helpers.RandomString(5))), "id of the peer for the network")
	noinput  = flag.Bool("noinput", false, "no input for automatic nodes")

	bossBool = flag.Bool("boss", false, "whether you want a boss or not")
	bossAddr = flag.String("bossAddr", "[::1]:3000", "boss location")
)

var self *Peer
var boss *Boss

var sendId string
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
	self.Address = fmt.Sprintf("[%s]:%s", myIp, *port)
}

func main() {

	fmt.Println("hey")
	incomingConnection := setupIncomingConnection(self.Address)

	if *bossBool {

		boss = SetupBossOnAddress(*bossAddr)
		go boss.ListenAndHandleBoss()
	}

	fmt.Println("ok")
	fmt.Println(self.Id, "listening on", self.Address, "scanning on", *scanPort)
	fmt.Println("Network: ", self)

	inputCb := make(chan string)
	connectionCb := make(chan []byte)
	searcherCb := make(chan *net.UDPConn)

	if *noinput == false {

		go runReadInput(inputCb)
	}
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

func runReadInput(cb chan string) {

	for {

		var input string
		n, err := fmt.Scanf("%s\n", &input)

		if n > 0 {

			panicOnError(err)
			cb <- input
		}
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

func inputHandler(input string) {

	switch messageState {
	case 0:

		str := string(input)

		if str == "connect" {

			messageState = CONNECTION_STATE

		} else {

			sendId = input
			messageState = MESSAGE_STATE

		}

	case 1:

		messageState = START_STATE

		message := &Message{Body: input}

		self.SendMessage(message, sendId)

	case 2:

		messageState = START_STATE
		self.StablishConnection(input)
	}
}

func panicOnError(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}
