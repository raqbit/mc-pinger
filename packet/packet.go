package packet

import (
	"bytes"
	enc "github.com/Raqbit/mc-pinger/encoding"
)

// Represents a Minecraft packet.
type Packet interface {
	PacketID() enc.VarInt
}

// Packet which is able to be encoded.
type EncodablePacket interface {
	Packet
	Marshal() ([]byte, error)
}

// Packet which is able to be decoded.
type DecodablePacket interface {
	Packet
	Unmarshal(reader enc.Reader) error
}

// Write a packet to the given Writer.
func WritePacket(p EncodablePacket, wr enc.Writer) error {
	// Marshal packet data
	data, err := p.Marshal()

	if err != nil {
		return err
	}

	// Get the packet id in packed form
	pId, err := getPacketIdBytes(p)

	if err != nil {
		return err
	}

	// Calculate packet length
	length := enc.VarInt(len(pId) + len(data))

	// Write packet length
	if err = enc.WriteVarInt(wr, length); err != nil {
		return err
	}

	// Write packet id
	if _, err = wr.Write(pId); err != nil {
		return err
	}

	// Write packet data
	if _, err = wr.Write(data); err != nil {
		return err
	}

	return nil
}

// Get the packet ID of given packet in byte form.
func getPacketIdBytes(p Packet) ([]byte, error) {
	packetId := p.PacketID()

	pIdBuff := new(bytes.Buffer)

	err := enc.WriteVarInt(pIdBuff, packetId)

	if err != nil {
		return nil, err
	}

	return pIdBuff.Bytes(), nil
}

// Reads a packet header (length, version) from the given Reader.
func ReadPacketHeader(rd enc.Reader) (enc.VarInt, enc.VarInt, error) {
	pLen, err := enc.ReadVarInt(rd)

	if err != nil {
		return 0, 0, err
	}

	pId, err := enc.ReadVarInt(rd)

	if err != nil {
		return 0, 0, err
	}

	return pLen, pId, nil
}
