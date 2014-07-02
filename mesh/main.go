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
	id              = flag.String("id", helpers.SHA1([]byte(helpers.RandomString(5))), "id of the peer for the network")
)

var self *Peer
var sentMessages []string

var currentMessage *Message
var messageState int = 0

func init() {

	flag.Parse()

	self = new(Peer)
	self.Id = *id
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
			if messageState == 0 {

				var dest_id = string(input)
				var next_peer = self.FindNearestPeerToId(dest_id)

				if next_peer != nil {

					currentMessage = &Message{Destination: next_peer.Address, FinalDestinationId: dest_id}
					messageState = 1

					fmt.Println("Sending message to", dest_id, "through", next_peer)
				} else {

					fmt.Println("Couldn't find peer with that id")
				}

			} else {

				messageState = 0
				currentMessage.Body = string(input)
				currentMessage.Origin = *self
				connection := setupOutgoingConnection(currentMessage.Destination)
				writeOutput(generateJSON(currentMessage), connection)
			}

		case connection := <-connectionCb:

			go incomingConnectionHandler(connection)

		case connection := <-searcherCb:

			go foundConnectionHandler(connection)
		}
	}
}

func runReadInput(cb chan []byte) {

	for {

		in := readInput(os.Stdin)
		cb <- in[:len(in)-1]
	}
}
func runConnectionInput(connection net.Listener, cb chan net.Conn) {
	for {

		con, err := connection.Accept()
		panicOnError(err)
		cb <- con
	}
}

func searchPeersOnPort(port string, cb chan net.Conn) {

	for {

		network := []string{"10.0.5.33"}
		for _, address := range network {

			var tcpAd = address + ":" + port

			tcpAddress, err := net.ResolveTCPAddr("tcp", tcpAd)
			panicOnError(err)

			//Check if is already a peer
			isPeer := false
			for _, p := range self.ConnectedPeers {
				if tcpAddress.String() == p.Address {
					isPeer = true
					break
				}
			}

			if tcpAddress.String() != self.Address && !isPeer {
				//Not looking for myself nor a peer already connected

				tcpConnection, err := net.DialTCP("tcp", nil, tcpAddress)

				if err == nil && tcpConnection != nil {

					cb <- tcpConnection
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func foundConnectionHandler(connection net.Conn) {

	fmt.Println("Found a connection opened. Sending my peer info...")
	mes := &Message{Origin: *self}
	mes.AssignRandomID()
	messageSent(mes.Id)
	writeOutput(generateJSON(mes), connection)
	connection.Close()
}

func incomingConnectionHandler(connection net.Conn) {

	message := parseJSON(readInput(connection))
	resp := isResponse(message.Id)

	if message.Body == "" {

		if err := self.AddConnectedPeer(message.Origin); err == nil {
			fmt.Println("Added peer: self ->", self)
		}

		if !resp {
			var respAddress = message.Origin.Address
			message.Origin = *self
			writeOutput(generateJSON(message), setupOutgoingConnection(respAddress))
		}
	} else {

		fmt.Println("! Message from ", message.Origin.Address, " -> ", message.Body)
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

func messageSent(id string) {

	sentMessages = append(sentMessages, id)
}

func isResponse(id string) bool {

	for _, m := range sentMessages {

		if id == m {

			return true
		}
	}

	return false
}
