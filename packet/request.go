package packet

import (
	enc "mc-pinger/encoding"
)

type RequestPacket struct{}

func (RequestPacket) PacketID() enc.VarInt {
	return 0x00
}

func (h RequestPacket) Marshal() ([]byte, error) {
	// Packet does not have any content.
	return make([]byte, 0), nil
}
