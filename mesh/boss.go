package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	InfoType = iota + 1
	MessageType
	ConnectType
)

type BossPacket struct {
	Type        int         `json:"type"`
	Data        string      `json:"data,omitempty"`
	PeerData    Peer        `json:"peerData,omitempty"`
	MessageData BossMessage `json:"messageData,omitempty"`
}

type BossMessage struct {
	From string `json:"from,omitempty"`
	To   string `json:"to"`
	Body string `json:"body,omitempty"`
}

type Boss struct {
	Connection *net.TCPConn
}

func SetupBossOnAddress(address string) *Boss {

	addr, err := net.ResolveTCPAddr("tcp6", address)
	panicOnError(err)

	con, err := net.DialTCP("tcp6", nil, addr)
	panicOnError(err)

	return &Boss{Connection: con}
}

func (b *Boss) ListenAndHandleBoss() {

	fmt.Println("Connected to boss", b.Connection.RemoteAddr())
	for {

		var buffer []byte = make([]byte, 4096)
		n, err := b.Connection.Read(buffer[0:])
		panicOnError(err)

		if n > 0 {

			packet := new(BossPacket)
			json.Unmarshal(buffer[:n], packet)

			switch packet.Type {

			case InfoType:
				go b.SendPeerInfo(self)

			case ConnectType:
				go self.StablishConnection(packet.Data)
			}
		}
	}
}

func (b *Boss) SendPeerInfo(p *Peer) {

	packet := new(BossPacket)

	packet.Type = InfoType
	packet.PeerData = *p

	pa, err := json.Marshal(packet)
	panicOnError(err)

	writeOutput(pa, b.Connection)
}
