package main

import (
	"net"
)

func main() {

	port := ":8080"
	addr, err := net.ResolveUDPAddr("udp", port)
	panicOnError(err)

	con, err := net.DialUDP("udp", nil, addr)
	panicOnError(err)

	b := []byte("Hello")
	_, err = con.Write(b)
	panicOnError(err)

}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
