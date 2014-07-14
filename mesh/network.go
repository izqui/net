package main

import (
	"errors"
	"net"
)

func setupIncomingConnection(address string) *net.UDPConn {

	addr, err := net.ResolveUDPAddr("udp", address)
	panicOnError(err)

	con, err := net.ListenUDP("udp", addr)
	panicOnError(err)

	return con
}

func setupOutgoingConnection(address string) *net.UDPConn {
	udpAddress, err := net.ResolveUDPAddr("udp", address)
	panicOnError(err)
	udpConnection, err := net.DialUDP("udp", nil, udpAddress)
	panicOnError(err)

	return udpConnection
}

func pingAddress(address string) (connection *net.UDPConn, err error) {

	udpAddress, err := net.ResolveUDPAddr("udp", address)
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

		return net.DialUDP("udp", nil, udpAddress)
	}

	return nil, errors.New("You are already connected")
}

func myIp() string {

	ifaces, err := net.Interfaces()
	if err != nil {
		panicOnError(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			panicOnError(err)
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	panicOnError(errors.New("are you connected to the network?"))
	return ""
}
