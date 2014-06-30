package main

import (
	"encoding/json"
	"flag"

	"io"
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

	input := readInput(os.Stdin)
	mes := &Message{Body: string(input)}
	writeJSON(mes, os.Stdout)
}

func readInput(reader io.Reader) []byte {

	/*stat, err := reader.Stat()
	panicOnError(err)
	size := stat.Size()*/

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

func panicOnError(err error) {

	if err != nil {

		panic(err)
	}
}
