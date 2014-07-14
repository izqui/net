package main

import (
	"fmt"
	"net"
)

func main() {

	a, _ := net.ResolveTCPAddr("tcp4", ":1080")
	l, _ := net.ListenTCP("tcp4", a)

	for {

		l.Accept()
		fmt.Println("Con")
	}
}
