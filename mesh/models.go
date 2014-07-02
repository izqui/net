package main

import (
	"errors"
	"fmt"

	"github.com/izqui/helpers"
)

type Peer struct {
	Id             string `json:"id"`
	Address        string `json:"address"`
	ConnectedPeers []Peer `json:"connected_peers,omitempty"`
}

type Message struct {
	Id            string `json:"id"`
	Body          string `json:"body,omitempty"`
	Origin        Peer   `json:"origin_peer"`
	DestinationId string `json:"destination_id"`
}

func (m *Message) AssignRandomID() {

	m.Id = helpers.SHA1([]byte(helpers.RandomString(10)))
}

func (p *Peer) AddConnectedPeer(newPeer Peer) error {

	if p.Id == newPeer.Id {

		return errors.New("You are trying to add yourself as a peer")
	}

	for _, con := range p.ConnectedPeers {

		if con.Id == newPeer.Id {

			return errors.New("Peer was already connected")
		}
	}

	p.ConnectedPeers = append(p.ConnectedPeers, newPeer)
	p.removeIfPresent(p.Id)

	return nil
}

func (p *Peer) removeIfPresent(id string) {

	connected := p.ConnectedPeers

	for i, c := range p.ConnectedPeers {

		if c.Id == id {

			connected = remove(connected, i)
		} else {

			fmt.Println("dont remove")
		}

		c.removeIfPresent(id)
	}

	p.ConnectedPeers = connected
}

func remove(slice []Peer, i int) []Peer {

	fmt.Println("remove")
	copy(slice[i:], slice[i+1:])
	slice[len(slice)-1] = Peer{}
	return slice[:len(slice)-1]
}
