package p2p

import (
	"log"
	"net"
	"sync"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder
	mu         sync.RWMutex
	peers      map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
		shakeHands: NOPHandshakeFunc,
	}
}

func (slf *TCPTransport) ListenAndAccept() error {
	var err error
	slf.listener, err = net.Listen("tcp", slf.listenAddr)
	if err != nil {
		return err
	}
	go slf.startAcceptLoop()
	return nil
}

func (slf *TCPTransport) startAcceptLoop() {
	for {
		conn, err := slf.listener.Accept()
		if err != nil {
			log.Printf("TCP accept error: %s\n", err)
			continue
		}
		log.Printf("new incoming connection: %+v\n", conn)
		go slf.handleConn(conn)
	}
}

type Temp struct{}

func (slf *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	if err := slf.shakeHands(peer); err != nil {
		log.Printf("handshake failed: %s\n", err)
		_ = conn.Close()
		return
	}
	msg := &Temp{}
	for {
		if err := slf.decoder.Decode(conn, msg); err != nil {
			log.Printf("TCP decode error: %s\n", err)
			continue
		}
	}
}
