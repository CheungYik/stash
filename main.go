package main

import (
	"log"
	"stash/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":8080")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal()
	}
	select {}
}
