package main

import (
	"encoding/json"
	"fmt"
	"github.com/izqui/functional"
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
	Connections []string
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

func (n *Node) ConnectToNode(addr string) {

	b, err := json.Marshal(BossPacket{Type: ConnectType, Data: addr})
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

				n.Id = packet.PeerData.Id
				n.PeerAddress = packet.PeerData.Address

				//Figure out what connections to do
				toadd := functional.Filter(func(p Peer) (f bool) {

					if p.Id != n.Id {

						for _, node := range n.Connections {

							if node == p.Id {
								//If we already added it, we are not interested
								return
							}
						}
						for _, node := range nodes {

							if node.Id == p.Id {
								//If it is among nodes in boss, add it
								f = true
							}
						}
					}

					return

				}, packet.PeerData.ConnectedPeers)

				n.Connections = append(n.Connections, functional.DoMap(func(p Peer) string { return p.Id }, toadd).([]string)...)
				fmt.Println(n)

			}
		}
	}
}
