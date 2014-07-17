package main

import (
	"errors"
	"net"
)

func setupIncomingConnection(address string) *net.UDPConn {

	addr, err := net.ResolveUDPAddr("udp6", address)
	panicOnError(err)

	con, err := net.ListenUDP("udp6", addr)
	panicOnError(err)

	return con
}

func setupOutgoingConnection(address string) *net.UDPConn {
	udpAddress, err := net.ResolveUDPAddr("udp6", address)
	panicOnError(err)
	udpConnection, err := net.DialUDP("udp6", nil, udpAddress)
	panicOnError(err)

	return udpConnection
}

func pingAddress(address string) (connection *net.UDPConn, err error) {

	udpAddress, err := net.ResolveUDPAddr("udp6", address)
	if err != nil {

		return nil, err
	}

	//Check if is already a peer
	isPeer := false
	for _, p := range self.ConnectedPeers {
		if udpAddress.String() == p.Address {
			isPeer = true
			break
		}
	}

	if udpAddress.String() != self.Address && !isPeer {
		//Not looking for myself nor a peer already connected

		return net.DialUDP("udp6", nil, udpAddress)
	}

	return nil, errors.New("You are already connected")
}
