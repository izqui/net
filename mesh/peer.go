package main

import (
	"errors"
	"fmt"
	"github.com/izqui/helpers"
	"net"
	"sort"
)

type PeerSlice []*Peer

func (slice PeerSlice) Len() int {

	return len(slice)
}

func (slice PeerSlice) Swap(i, j int) {

	slice[i], slice[j] = slice[j], slice[i]
}

func (slice PeerSlice) Less(i, j int) bool {

	return slice[i].Id < slice[j].Id
}

func (slice PeerSlice) remove(i int) PeerSlice {

	copy(slice[i:], slice[i+1:])
	slice[len(slice)-1] = nil
	return slice[:len(slice)-1]
}

type Peer struct {
	Id             string    `json:"id"`
	Address        string    `json:"address"`
	ConnectedPeers PeerSlice `json:"connected_peers,omitempty"`
	messagesSent   []string  `json:"-"`
}

func (p *Peer) String() string {

	con := ""
	if len(p.ConnectedPeers) > 0 {

		for _, c := range p.ConnectedPeers {

			if c != nil {
				con += c.String()
			}
		}

		con = fmt.Sprintf(" -> [%s]", con)
	}

	return fmt.Sprintf("%s%s", p.Id, con)
}

func (p *Peer) AddConnectedPeer(newPeer *Peer) error {

	if p.Id == newPeer.Id {

		return errors.New("You are trying to add yourself as a peer")
	}

	newPeer.removeIfPresent(p.Id)

	location := -1

	for i, con := range p.ConnectedPeers {

		if con.Id == newPeer.Id {

			if con.Hash() == newPeer.Hash() {

				return errors.New("Peer was already connected")
			} else {

				// Insert at this location, same id, different shit inside
				location = i
			}
		}
	}

	if location == -1 {

		p.ConnectedPeers = append(p.ConnectedPeers, newPeer)
	} else {
		p.ConnectedPeers[location] = newPeer
	}

	//Remove myself if I'm referenced in other peers
	//p.removeIfPresent(p.Id)

	return nil
}

func (p *Peer) Hash() string {

	data := p.Id

	if len(p.ConnectedPeers) > 0 {

		sort.Sort(p.ConnectedPeers)

		for _, c := range p.ConnectedPeers {

			if c != nil {

				data += c.Hash()
			}
		}
	}

	return helpers.SHA1([]byte(data))
}

func (p *Peer) removeIfPresent(id string) {

	count := len(p.ConnectedPeers)
	i := 0

	for i < count {

		c := p.ConnectedPeers[i]

		if c != nil {

			if c.Id == id {

				p.ConnectedPeers = p.ConnectedPeers.remove(i)
				count -= 1
				i -= 1

			} else {

				c.removeIfPresent(id)
			}
		}
		i++
	}
}

func (p *Peer) FindNearestPeerToId(id string) *Peer {

	for _, c := range p.ConnectedPeers {

		if c.Id == id {

			return c
		}
	}

	distance := 1000
	var peer *Peer = nil

	for _, c := range p.ConnectedPeers {

		n := c.distanceToId(id)

		if n > -1 && n < distance {

			peer = c
			distance = n
		}
	}

	return peer
}

func (p Peer) distanceToId(id string) int {

	for _, c := range p.ConnectedPeers {

		if c != nil && c.Id == id {

			return 1
		}

		n := c.distanceToId(id)
		if n > -1 {

			return 1 + n
		}
	}

	return -1
}

func (p *Peer) HandleIncomingConnection(input []byte) {

	message := parseJSON(input)
	resp := p.isExistingMessage(message.Id)

	if message.Body == "" {

		if err := self.AddConnectedPeer(message.Origin); err == nil {

			fmt.Println("Network: ", self)
			for _, p := range self.ConnectedPeers {

				if (resp && p.Id != message.Origin.Id) || !resp {

					respMessage := &Message{Id: message.Id, Origin: self}

					if boss != nil {
						boss.SendPeerInfo(self)
					}

					p.send(respMessage, p.Address)
				}
			}
		}

	} else if message.FinalDestinationId == self.Id {
		//Message is for me :)
		boss.SendMessageFlowInfo(message.Origin.Id, self.Id)
		fmt.Println("! Message from ", message.Origin.Id, " -> ", message.Body)
	} else {
		//Message is not for me. Redirecting it.
		fmt.Println("Broadcasting message from", message.Origin.Id, "to", message.FinalDestinationId)

		var next_peer = p.FindNearestPeerToId(message.FinalDestinationId)

		if next_peer != nil {

			message.Destination = next_peer.Address

			boss.SendMessageFlowInfo(self.Id, next_peer.Id)
			p.send(message, message.Destination)

			fmt.Println("Sending message to", message.FinalDestinationId, "through", next_peer)

		} else {

			fmt.Println("Couldn't find peer with that id")
		}
	}
}

func (p *Peer) HandleConnectionFound(connection *net.UDPConn) {

	mes := &Message{Origin: self}
	mes.AssignRandomID()

	if !p.isExistingMessage(mes.Id) {
		p.messagesSent = append(p.messagesSent, mes.Id)
	}

	writeOutput(generateJSON(mes), connection)
}

func (p *Peer) isExistingMessage(id string) bool {

	for _, m := range p.messagesSent {

		if id == m {

			return true
		}
	}

	return false
}

func (p *Peer) SendMessage(mes *Message, toPeerId string) {

	var next_peer = self.FindNearestPeerToId(toPeerId)

	if next_peer != nil {

		mes.Destination = next_peer.Address
		mes.FinalDestinationId = toPeerId
		messageState = MESSAGE_STATE

		fmt.Println("Sending message to", toPeerId, "through", next_peer.Id)

	} else {

		fmt.Println("Couldn't find peer with that id")
		return
	}

	mes.Origin = p
	mes.AssignRandomID()

	p.send(mes, next_peer.Address)
}

func (p *Peer) send(message *Message, address string) {

	connection := setupOutgoingConnection(address)
	writeOutput(generateJSON(message), connection)

	if !p.isExistingMessage(message.Id) {
		p.messagesSent = append(p.messagesSent, message.Id)
	}
}

func (p *Peer) StablishConnection(address string) {

	con, err := pingAddress(address)
	if err == nil && con != nil {

		mes := &Message{Origin: self}
		mes.AssignRandomID()
		p.send(mes, address)

	} else {
		fmt.Println("Couldn't stablish connection")
	}
}
