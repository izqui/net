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
var messageState int = 0

func init() {

	flag.Parse()

	self = new(Peer)
	self.Id = *id
	self.Address = fmt.Sprintf("%s:%s", myIp(), *port)
}

func main() {

	incomingConnection := setupIncomingConnection(self.Address)

	fmt.Println(self.Id, "listening on", self.Address, "scanning on", *scanPort)

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
			go incomingConnectionHandler(input)

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
func runConnectionInput(connection *net.UDPConn, cb chan []byte) {

	for {

		var buffer []byte = make([]byte, 512)
		n, addr, err := connection.ReadFromUDP(buffer[0:])
		panicOnError(err)

		if addr != nil {

			fmt.Println(addr)
			fmt.Println(n)
			cb <- buffer[:n]
		}
	}
}

func searchPeersOnPort(port string, cb chan *net.UDPConn) {

	for {

		network := []string{scanAddr}
		for _, address := range network {

			var tcpAd = address + ":" + port

			udpAddress, err := net.ResolveUDPAddr("udp", tcpAd)
			panicOnError(err)

			//Check if is already a peer
			isPeer := false
			for _, p := range self.ConnectedPeers {
				if udpAddress.String() == p.Address {
					isPeer = true
					break
				}
			}

			if udpAddress.String() != self.Address && !isPeer {
				//Not looking for myself nor a peer already connected

				udpConnection, err := net.DialUDP("udp", nil, udpAddress)

				if err == nil && udpConnection != nil {

					cb <- udpConnection
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func inputHandler(input []byte) {

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
		currentMessage.Origin = self
		currentMessage.AssignRandomID()
		messageSent(currentMessage.Id)

		connection := setupOutgoingConnection(currentMessage.Destination)
		writeOutput(generateJSON(currentMessage), connection)
	}
}
func foundConnectionHandler(connection *net.UDPConn) {

	fmt.Println("Found a connection opened. Sending my peer info...")
	mes := &Message{Origin: self}
	mes.AssignRandomID()
	messageSent(mes.Id)
	writeOutput(generateJSON(mes), connection)
	//connection.Close()
}

func incomingConnectionHandler(input []byte) {

	message := parseJSON(input)
	resp := isResponse(message.Id)

	if message.Body == "" {

		if err := self.AddConnectedPeer(message.Origin); err == nil {
			fmt.Println("Self ", self)
		} else {

			fmt.Println(err)
		}

		if !resp {

			messageSent(message.Id)
			message.Origin = self
			//Broadcast message
			for _, p := range self.ConnectedPeers {

				writeOutput(generateJSON(message), setupOutgoingConnection(p.Address))
			}
		}

	} else if message.FinalDestinationId == self.Id {
		//Message is for me :)
		fmt.Println("! Message from ", message.Origin.Id, " -> ", message.Body)
	} else {
		//Message is not for me. Redirecting it.
		fmt.Println("Broadcasting message from", message.Origin.Id, "to", message.FinalDestinationId)
		var next_peer = self.FindNearestPeerToId(message.FinalDestinationId)

		if next_peer != nil {

			message.Destination = next_peer.Address
			connection := setupOutgoingConnection(message.Destination)
			writeOutput(generateJSON(message), connection)

			fmt.Println("Sending message to", message.FinalDestinationId, "through", next_peer)

		} else {

			fmt.Println("Couldn't find peer with that id")
		}
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

func setupIncomingConnection(address string) *net.UDPConn {

	addr, err := net.ResolveUDPAddr("udp", address)
	panicOnError(err)

	con, err := net.ListenUDP("udp", addr)
	panicOnError(err)

	return con
}

func setupOutgoingConnection(address string) *net.UDPConn {
	udpAddress, err := net.ResolveUDPAddr("udp", address)
	panicOnError(err)
	udpConnection, err := net.DialUDP("udp", nil, udpAddress)
	panicOnError(err)

	return udpConnection
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
