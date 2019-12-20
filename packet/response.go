package packet

import (
	enc "github.com/Raqbit/mc-pinger/encoding"
	"io"
)

type ResponsePacket struct {
	Json enc.String
}

func (ResponsePacket) ID() enc.VarInt {
	return 0x00
}

func (rp *ResponsePacket) Unmarshal(reader io.Reader) error {
	// Read JSON string
	str, err := enc.ReadString(reader)

	if err != nil {
		return err
	}

	rp.Json = str

	return nil
}
