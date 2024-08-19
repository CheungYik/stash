package p2p

import (
	"errors"
	"io"
	"log"
	"net"
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

func (slf *TCPPeer) Close() error {
	return slf.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr string
	Handshake  HandshakeFunc
	Decoder    Decoder
	OnPeer     func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcCh    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

func (slf *TCPTransport) Consume() <-chan RPC {
	return slf.rpcCh
}

func (slf *TCPTransport) ListenAndAccept() error {
	var err error
	slf.listener, err = net.Listen("tcp", slf.ListenAddr)
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

func (slf *TCPTransport) handleConn(conn net.Conn) {
	var err error
	peer := NewTCPPeer(conn, true)
	defer func() {
		log.Printf("dropping peer connection: %s\n", err)
		peer.Close()
	}()
	if err = slf.Handshake(peer); err != nil {
		return
	}
	if slf.OnPeer != nil {
		if err = slf.OnPeer(peer); err != nil {
			return
		}
	}
	rpc := RPC{From: conn.RemoteAddr()}
	for {
		if err = slf.Decoder.Decode(conn, &rpc); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				break
			}
			log.Printf("TCP decode error: %s\n", err)
			continue
		}
		slf.rpcCh <- rpc
	}
}
