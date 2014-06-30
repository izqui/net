package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {

	//lookupIP()
	requestHead(lookupTCP())
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

func requestHead(tcpAddr *net.TCPAddr) {

	fmt.Printf("Requesting HEAD for %s:%d\n", tcpAddr.IP, tcpAddr.Port)
	connection, err := net.DialTCP("tcp", nil, tcpAddr)
	errorHandling(err)

	_, err = connection.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	errorHandling(err)

	result, err := ioutil.ReadAll(connection)
	errorHandling(err)

	fmt.Println(string(result))

}

func errorHandling(err error) {

	if err != nil {

		panic(err)
	}
}
