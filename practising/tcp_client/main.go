package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	//lookupIP()
	sendData(lookupTCP())
}

func lookupIP() {

	args := os.Args
	if len(args) != 2 {

		panic("Incorrect usage of program")
	} else {

		name := args[1]
		ips, err := net.LookupIP(name)

		if err == nil {

			for _, ip := range ips {

				fmt.Println(ip)
			}
		}

	}
}

func lookupTCP() *net.TCPAddr {

	args := os.Args
	if len(args) != 2 {

		panic("Incorrect usage of program")
	}

	name := args[1]

	tcp, err := net.ResolveTCPAddr("tcp", name)

	if err != nil {

		panic(err)
	}

	return tcp
}

func sendData(tcpAddr *net.TCPAddr) {

	fmt.Printf("Sending data to %s:%d\n", tcpAddr.IP, tcpAddr.Port)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	errorHandling(err)

	_, err = connection.Write([]byte("bullshit"))
	errorHandling(err)

	var buf [512]byte
	n, err := connection.Read(buf[0:])
	errorHandling(err)

	fmt.Println(string(buf[0:n]))

}

func errorHandling(err error) {

	if err != nil {

		panic(err)
	}
}
