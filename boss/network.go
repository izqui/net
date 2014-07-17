package main

import (
	"net"
)

type ConnectionCallback chan *net.TCPConn

func setupTCPListener(port string) *net.TCPListener {

	addr, err := net.ResolveTCPAddr("tcp6", "[::1]:"+port)
	panicOnError(err)

	listener, err := net.ListenTCP("tcp6", addr)
	panicOnError(err)

	return listener
}

func listenTCP(listener *net.TCPListener, cb ConnectionCallback) {

	for {

		conn, err := listener.AcceptTCP()
		panicOnError(err)

		cb <- conn
	}
}
