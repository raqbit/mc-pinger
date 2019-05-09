package encoding

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const (
	VarIntMaxByteSize = 5
)

var (
	// ErrVarIntTooLarge is returned when a read varint was too large
	// (more than 5 bytes)
	ErrVarIntTooLarge = errors.New("VarInt too large")
)

// Minecraft Protocol UnsignedShort type
type UnsignedShort uint16

// WriteUnsignedShort writes the passed UnsignedShort to the writer
func WriteUnsignedShort(buff io.Writer, value UnsignedShort) error {
	return binary.Write(buff, binary.BigEndian, uint16(value))
}

//ReadUnsignedShort reads an UnsignedShort from the reader
func ReadUnsignedShort(buff io.Reader) (UnsignedShort, error) {
	var short uint16
	err := binary.Read(buff, binary.BigEndian, &short)
	return UnsignedShort(short), err
}

// Minecrat Protocol UnsignedByte type
type UnsignedByte byte

func ReadUnsignedByte(r io.Reader) (UnsignedByte, error) {
	var bytes [1]byte
	_, err := r.Read(bytes[:1])
	return UnsignedByte(bytes[0]), err
}

func WriteUnsignedByte(w io.Writer, value UnsignedByte) error {
	var bytes [1]byte
	bytes[0] = byte(value)
	_, err := w.Write(bytes[:1])
	return err
}

// Minecraft Protocol VarInt type
type VarInt int32

// WriteVarInt writes the passed VarInt encoded integer to the writer.
func WriteVarInt(w io.Writer, value VarInt) error {
	for cont := true; cont; cont = value != 0 {
		temp := byte(value & 0x7F)

		// Casting value to a uint to get a logical shift
		value = VarInt(uint32(value) >> 7)

		if value != 0 {
			temp |= 0x80
		}

		if err := WriteUnsignedByte(w, UnsignedByte(temp)); err != nil {
			return err
		}
	}

	return nil
}

// ReadVarInt reads a VarInt encoded integer from the reader.
func ReadVarInt(r io.Reader) (VarInt, error) {
	var numRead uint
	var result int32
	var read UnsignedByte

	for cont := true; cont; cont = (read & 0x80) != 0 {
		var err error
		read, err = ReadUnsignedByte(r)

		if err != nil {
			return 0, err
		}

		value := read & 0x7F

		result |= int32(value) << (7 * numRead)

		numRead++

		if numRead > VarIntMaxByteSize {
			return 0, ErrVarIntTooLarge
		}
	}

	return VarInt(result), nil
}

type String string

// WriteString writes a VarInt prefixed utf-8 string to the
// writer.
func WriteString(w io.Writer, str String) error {

	// Creating buffer from string
	b := []byte(str)

	// Writing string length as varint to output buffer
	err := WriteVarInt(w, VarInt(len(b)))

	if err != nil {
		return err
	}

	// Writing string to buffer
	_, err = w.Write(b)

	return err
}

// ReadString reads a VarInt prefixed utf-8 string to the
// reader. It uses io.ReadFull to ensure all bytes are read.
func ReadString(r io.Reader) (String, error) {

	// Reading string size encoded as VarInt
	l, err := ReadVarInt(r)

	if err != nil {
		return "", nil
	}

	// Checking if string size is valid
	if l < 0 || l > math.MaxInt16 {
		return "", errors.New("string length out of bounds")
	}

	// Creating string buffer with the specified size
	stringBuff := make([]byte, int(l))

	// Reading l amount of bytes from the buffer
	_, err = io.ReadFull(r, stringBuff)

	return String(stringBuff), err
}
