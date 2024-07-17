package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":8080"
	tr := NewTCPTransport(listenAddr)
	assert.Equal(t, listenAddr, tr.listenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
