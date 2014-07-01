package main

import (
	"errors"
)

type Peer struct {
	Id             string `json:"id"`
	Address        string `json:"address"`
	ConnectedPeers []Peer `json:"connected_peers,omitempty"`
}

type Message struct {
	Body          string `json:"body,omitempty"`
	Origin        Peer   `json:"origin_peer"`
	DestinationId string `json:"destination_id"`
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
	return nil
}
