package main

import (
	"encoding/json"
	"io"
)

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
