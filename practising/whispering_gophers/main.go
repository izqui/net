package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"bufio"
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
	inputCb := make(chan string)
	connectionCb := make(chan net.Conn, 1)
	go runReadInput(inputCb)
	go runConnectionInput(incomingConnection, connectionCb)
	os.Stdout.Write([]byte("Message: "))
	for {
		select {
		case connection := <-connectionCb:
			fmt.Println("never being called")
			fmt.Println("data " + string(readInput(connection)))
			//message := parseJSON()
			//connection.Close()
			//writeOutput([]byte(message.Body), os.Stdout)
		case input := <-inputCb:
			mes := &Message{Body: input}
			connection := setupOutgoingConnection(*address)
			fmt.Println("[debug] sending out data")
			writeOutput(generateJSON(mes), connection)
		}
	}
}
func runReadInput(cb chan string) {
	in := bufio.NewReader(os.Stdin)
	for {
		l, err := in.ReadString('\n')
		panicOnError(err)
		fmt.Println("[debug] studin input")
		cb <- l
	}
}
func runConnectionInput(connection net.Listener, cb chan net.Conn) {
	for {
		fmt.Println("[debug] tcp checker")
		con, err := connection.Accept()
		fmt.Println("[debug] tcp2")
		panicOnError(err)
		fmt.Println("test")
		cb <- con
		fmt.Println("FJDSKLJ")
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
func parseJSON(data []byte) (message *Message) {
	err := json.Unmarshal(data, message)
	panicOnError(err)
	return
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
	if err != nil {
		panic(err)
	}
}
