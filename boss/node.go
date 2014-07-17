package main

import (
	"encoding/json"
	"fmt"
	"github.com/izqui/helpers"
	"net"
	"os/exec"
	"strconv"
)

const (
	InfoType = iota + 1
	MessageType
	ConnectType
)

type BossPacket struct {
	Type     int    `json:"type"`
	Data     string `json:"data, omitempty"`
	PeerData Peer   `json:"peerData, omitempty"`
}

type Peer struct {
	Id             string   `json:"id"`
	Address        string   `json:"address"`
	ConnectedPeers []Peer   `json:"connected_peers,omitempty"`
	messagesSent   []string `json:"-"`
}

type Node struct {
	BossConnection *net.TCPConn

	//Mesh
	PeerAddress string
	Id          string
}

func BootUpNode(id string, port int) {

	if port == 0 {

		port = helpers.RandomInt(8000, 9999)
	}

	cmd := exec.Command("mesh", "--id", id, "--port", strconv.Itoa(port), "--boss", "--noinput")
	cmd.Start()
	go cmd.Wait()
}

func (n *Node) GetInfo() {

	b, err := json.Marshal(BossPacket{Type: InfoType})
	panicOnError(err)

	_, err = n.BossConnection.Write(b)
	panicOnError(err)
}

func (n *Node) ListenForConnections() {

	for {

		var buffer []byte = make([]byte, 4096)
		nu, err := n.BossConnection.Read(buffer[0:])
		panicOnError(err)

		if nu > 0 {

			packet := new(BossPacket)
			json.Unmarshal(buffer[:nu], packet)

			switch packet.Type {

			case InfoType:

				fmt.Println(packet.Data, packet.PeerData)
				n.GetInfo()
			}
		}
	}
}
