package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Peer struct {
}
type Message struct {
	Body               string `json:"body"`
	OriginAddress      string `json:"origin_address"`
	DestinationAddress string `json:"destination_address"`
}

var (
	outgoingAddress = flag.String("out", "localhost:3003", "address of peer")
	port            = flag.String("port", "0", "your local port")
)

func init() {

	flag.Parse()
}

func main() {

	incomingConnection := setupIncomingConnection(myIp() + ":" + *port)

	fmt.Println("Listening on ", incomingConnection.Addr())

	inputCb := make(chan []byte)
	connectionCb := make(chan net.Conn)
	go runReadInput(inputCb)
	go runConnectionInput(incomingConnection, connectionCb)

	mes := new(Message)

	for {
		select {

		case input := <-inputCb:
			mes := &Message{Body: string(input), OriginAddress: incomingConnection.Addr().String()}
			connection := setupOutgoingConnection(*outgoingAddress)
			writeOutput(generateJSON(mes), connection)

		case connection := <-connectionCb:

			message := parseJSON(readInput(connection))
			fmt.Println("! Message from ", message.OriginAddress, " -> ", message.Body)
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

		con, err := connection.Accept()
		panicOnError(err)
		cb <- con
	}
}
func readInput(reader io.Reader) []byte {
	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	panicOnError(err)
	return buf[0:n]
}
func writeOutput(content []byte, writer io.Writer) {
	writer.Write(content)
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
