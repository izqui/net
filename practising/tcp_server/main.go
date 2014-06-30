package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	var args = os.Args
	if len(args) < 2 {

		panic("Incorrect use of the program")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%s", os.Args[1]))
	errorHandling(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	errorHandling(err)

	fmt.Printf("Listening TCP on %s:%d\n", tcpAddr.IP, tcpAddr.Port)
	for {

		connection, err := listener.Accept()
		errorHandling(err)

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {

	var buf [512]byte

	fmt.Printf("Opened TCP connection with %s\n", connection.RemoteAddr())
	for {

		n, err := connection.Read(buf[0:])
		if err != nil {

			fmt.Println(err)
			return
		}
		fmt.Println("Read data")
		_, err = connection.Write(buf[0:n])
		if err != nil {

			fmt.Println(err)
			return
		}
		fmt.Println("Send data")
	}
}
func errorHandling(err error) {

	if err != nil {

		panic(err)
	}
}
