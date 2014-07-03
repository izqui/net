package main

import "fmt"

func test() {

	p := &Peer{Id: "A"}
	p.AddConnectedPeer(&Peer{Id: "B"})
	fmt.Println(p)
	p.AddConnectedPeer(&Peer{Id: "B", ConnectedPeers: PeerSlice{&Peer{Id: "A"}, &Peer{Id: "A"}}})
	fmt.Println(p)
}

/*
func main(){

	test()
}*/
