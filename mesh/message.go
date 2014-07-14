package main

import (
	"github.com/izqui/helpers"
)

type Message struct {
	Id                 string `json:"id"`
	Body               string `json:"body,omitempty"`
	Origin             *Peer  `json:"origin_peer"`
	Destination        string `json:"-"`
	FinalDestinationId string `json:"destination_id"`
}

func (m *Message) AssignRandomID() {

	m.Id = helpers.SHA1([]byte(helpers.RandomString(10)))
}
