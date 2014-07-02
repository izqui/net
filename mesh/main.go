package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/izqui/helpers"
	"io"
	"net"
	"os"
	"time"
)

var (
	defAd   = "10.0.5.33"
	defPort = "3003"
)

var (
	outgoingAddress = flag.String("out", defAd+":"+defPort, "address of peer")
	port            = flag.String("port", "0", "your local port")
	name            = flag.String("name", helpers.SHA1([]byte(helpers.RandomString(5))), "name of the peer for the network")
)

var self *Peer

func init() {

	flag.Parse()

	self = new(Peer)
	self.Id = *name
	self.Address = fmt.Sprintf("%s:%s", myIp(), *port)

}

func main() {

	incomingConnection := setupIncomingConnection(self.Address)
	self.Address = incomingConnection.Addr().String()

	fmt.Println(self.Id, "listening on", self.Address)

	inputCb := make(chan []byte)
	connectionCb := make(chan net.Conn)
	searcherCb := make(chan net.Conn)

	go runReadInput(inputCb)
	go runConnectionInput(incomingConnection, connectionCb)
	go searchPeersOnPort(defPort, searcherCb)

	for {
		select {

		case input := <-inputCb:
			mes := &Message{Body: string(input), Origin: *self}
			connection := setupOutgoingConnection(*outgoingAddress)
			writeOutput(generateJSON(mes), connection)

		case connection := <-connectionCb:

			message := parseJSON(readInput(connection))

			if message.Body == "" {

				fmt.Println("add peer")
				self.AddConnectedPeer(message.Origin)
				fmt.Println(self)

			} else {

				fmt.Println("! Message from ", message.Origin.Address, " -> ", message.Body)
			}

		case connection := <-searcherCb:

			fmt.Println("Found a connection opened. Sending my peer info...")
			mes := &Message{Origin: *self}
			writeOutput(generateJSON(mes), connection)
		}
	}
}
func runReadInput(cb chan []byte) {

	for {

		cb <- readInput(os.Stdin)
	}
}
func runConnectionInput(connection net.Listener, cb chan net.Conn) {
	for {

		fmt.Println("waiting for connction")
		con, err := connection.Accept()
		fmt.Println("incoming connection from", con.RemoteAddr())
		panicOnError(err)
		cb <- con
	}
}

func searchPeersOnPort(port string, cb chan net.Conn) {

	for {

		fmt.Println("Checking for peers")
		network := []string{"10.0.5.33"}
		for _, address := range network {

			var tcpAd = address + ":" + port

			tcpAddress, err := net.ResolveTCPAddr("tcp", tcpAd)
			panicOnError(err)

			if tcpAddress.String() != self.Address {
				//Not looking for myself

				tcpConnection, err := net.DialTCP("tcp", nil, tcpAddress)

				if err == nil && tcpConnection != nil {

					cb <- tcpConnection
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func readInput(reader io.Reader) []byte {
	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	panicOnError(err)
	return buf[0:n]
}
func writeOutput(content []byte, writer io.Writer) {
	_, err := writer.Write(content)
	panicOnError(err)
}
func parseJSON(data []byte) *Message {

	message := new(Message)
	err := json.Unmarshal(data, message)
	panicOnError(err)
	return message
}
func generateJSON(mes *Message) []byte {
	data, err := json.Marshal(mes)
	panicOnError(err)
	return data
}
func setupIncomingConnection(address string) net.Listener {

	listener, err := net.Listen("tcp4", address)
	panicOnError(err)

	return listener
}
func setupOutgoingConnection(address string) net.Conn {
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	panicOnError(err)
	tcpConnection, err := net.DialTCP("tcp", nil, tcpAddress)
	panicOnError(err)
	return tcpConnection
}
func panicOnError(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func myIp() string {

	ifaces, err := net.Interfaces()
	if err != nil {
		panicOnError(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			panicOnError(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	panicOnError(errors.New("are you connected to the network?"))
	return ""
}
