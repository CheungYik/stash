package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8080"
	tcpOpts := TCPTransportOpts{
		ListenAddr: ":8080",
		Handshake:  NOPHandshakeFunc,
		Decoder:    DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)
	assert.Equal(t, listenAddr, tr.ListenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
