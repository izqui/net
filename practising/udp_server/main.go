package main

import (
	"fmt"
	"net"
)

func main() {

	cb := make(chan []byte)

	go setupServer(":8080", cb)

	for {
		select {

		case data := <-cb:
			fmt.Println("Data ", string(data))
		}
	}
}

func setupServer(port string, cb chan []byte) {

	addr, err := net.ResolveUDPAddr("udp", port)
	panicOnError(err)

	con, err := net.ListenUDP("udp", addr)
	fmt.Println("Listening", addr)
	panicOnError(err)

	for {

		var buffer []byte = make([]byte, 512)
		n, addr, err := con.ReadFromUDP(buffer[0:])
		panicOnError(err)

		if addr != nil {

			fmt.Println(addr)
			fmt.Println(n)
			cb <- buffer[:n]
		}
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
