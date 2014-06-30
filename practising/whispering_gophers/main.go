package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	address = flag.String("address", "localhost:3003", "address of peer")
)

type Message struct {
	Body string `json:"body"`
}

func init() {

	flag.Parse()
}

func main() {

	os.Stdout.Write([]byte("Message: "))

	input := readInput(os.Stdin)
	os.Stdout.Write([]byte(input))

	mes := &Message{Body: string(input)}

	connection := setupConnection(*address)
	writeJSON(mes, connection)
	fmt.Printf("Sent to %s\n", connection.RemoteAddr())

	_, err := os.Stdout.Write(readInput(connection))
	panicOnError(err)

	connection.Close()

}

func readInput(reader io.Reader) []byte {

	buf := make([]byte, 512)
	n, err := reader.Read(buf)
	panicOnError(err)

	return buf[0:n]
}

func writeJSON(mes *Message, writer io.Writer) {

	enc := json.NewEncoder(writer)
	err := enc.Encode(mes)
	panicOnError(err)
}

func setupConnection(address string) net.Conn {

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
