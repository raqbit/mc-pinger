package packet

import (
	enc "mc-pinger/encoding"
)

type ResponsePacket struct {
	Json enc.String
}

func (ResponsePacket) PacketID() enc.VarInt {
	return 0x00
}

func (rp *ResponsePacket) Unmarshal(reader enc.Reader) error {
	// Read JSON string
	str, err := enc.ReadString(reader)

	if err != nil {
		return err
	}

	rp.Json = str

	return nil
}
