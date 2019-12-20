package packet

import (
	enc "github.com/Raqbit/mc-pinger/encoding"
)

type RequestPacket struct{}

func (RequestPacket) ID() enc.VarInt {
	return 0x00
}

func (h RequestPacket) Marshal() ([]byte, error) {
	// Packet does not have any content.
	return make([]byte, 0), nil
}
