package packet

import (
	"bytes"
	enc "github.com/Raqbit/mc-pinger/encoding"
	"io"
)

type HandshakePacket struct {
	ProtoVer   enc.VarInt
	ServerAddr enc.String
	ServerPort enc.UnsignedShort
	NextState  enc.VarInt
}

func (h HandshakePacket) Marshal() ([]byte, error) {
	var buffer bytes.Buffer
	err := h.write(&buffer)
	return buffer.Bytes(), err
}

func (HandshakePacket) ID() enc.VarInt {
	return 0x00
}

func (h HandshakePacket) write(buffer io.Writer) error {

	// Write protocol version
	if err := enc.WriteVarInt(buffer, h.ProtoVer); err != nil {
		return err
	}

	// Write server address
	if err := enc.WriteString(buffer, h.ServerAddr); err != nil {
		return err
	}

	// Write server port
	if err := enc.WriteUnsignedShort(buffer, h.ServerPort); err != nil {
		return err
	}

	// Write next connection state
	if err := enc.WriteVarInt(buffer, h.NextState); err != nil {
		return err
	}

	return nil
}
