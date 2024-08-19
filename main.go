package main

import (
	"log"
	"stash/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr: ":8080",
		Handshake:  p2p.NOPHandshakeFunc,
		Decoder:    p2p.DefaultDecoder{},
		OnPeer:     OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	go func() {
		for {
			msg := <-tr.Consume()
			log.Printf("Received message: %+v", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal()
	}
	select {}
}
