package main

import (
	"encoding/json"
	"flag"
	"io"
	"net"
	"os"
)

var (
	address = flag.String("address", "localhost:3003", "address of peer")
	port    = flag.String("port", "8080", "your port")
)

type Message struct {
	Body string `json:"body"`
	//Address string `json:"address"`
}

func init() {

	flag.Parse()
}

func main() {
	incomingConnection := setupIncomingConnection(*port)
	inputCb := make(chan []byte)
	connectionCb := make(chan net.Conn)
	go runReadInput(inputCb)
	go runConnectionInput(incomingConnection, connectionCb)
	os.Stdout.Write([]byte("Message: "))

	for {
		select {

		case input := <-inputCb:
			mes := &Message{Body: string(input)}
			connection := setupOutgoingConnection(*address)
			writeOutput(generateJSON(mes), connection)

		case connection := <-connectionCb:

			message := parseJSON(readInput(connection))
			writeOutput([]byte(message.Body), os.Stdout)
		}
	}
}
func runReadInput(cb chan []byte) {

	for {

		input := readInput(os.Stdin)
		cb <- input
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
func setupIncomingConnection(port string) net.Listener {
	tcpAddress, err := net.ResolveTCPAddr("tcp", string(":"+port))
	panicOnError(err)
	listener, err := net.ListenTCP("tcp", tcpAddress)
	panicOnError(err)
	writeOutput([]byte("Accepting incoming connections on port "+string(port)+" \n"), os.Stdout)
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
		println(err == io.EOF)
	}
}
