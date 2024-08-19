package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (slf GOBDecoder) Decode(reader io.Reader, msg *RPC) error {
	return gob.NewDecoder(reader).Decode(msg)
}

type DefaultDecoder struct{}

func (slf DefaultDecoder) Decode(reader io.Reader, msg *RPC) error {
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	return nil
}
