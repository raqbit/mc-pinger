package encoding

import (
	"bytes"
	"io"
	"math"
	"testing"
)

func TestWriteUnsignedShort(t *testing.T) {
	tests := []struct {
		Value    UnsignedShort
		Expected []byte
	}{
		{Value: 0, Expected: []byte{0x00, 0x00}},
		{Value: 1, Expected: []byte{0x00, 0x01}},
		{Value: 2, Expected: []byte{0x00, 0x02}},
		{Value: 127, Expected: []byte{0x00, 0x7f}},
		{Value: 128, Expected: []byte{0x00, 0x80}},
		{Value: 255, Expected: []byte{0x00, 0xff}},
		{Value: math.MaxUint16, Expected: []byte{0xff, 0xff}},
	}

	var buff bytes.Buffer
	_ = io.Writer(&buff)

	for _, test := range tests {
		err := WriteUnsignedShort(&buff, test.Value)

		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(test.Expected, buff.Bytes()) != 0 {
			// Not equal
			t.Errorf("Unable to convert %d: %v != %v", test.Value, buff.Bytes(), test.Expected)
		}

		buff.Reset()
	}
}

func TestReadUnsignedShort(t *testing.T) {
	tests := []struct {
		Expected UnsignedShort
		Value    []byte
	}{
		{Expected: 0, Value: []byte{0x00, 0x00}},
		{Expected: 1, Value: []byte{0x00, 0x01}},
		{Expected: 2, Value: []byte{0x00, 0x02}},
		{Expected: 127, Value: []byte{0x00, 0x7f}},
		{Expected: 128, Value: []byte{0x00, 0x80}},
		{Expected: 255, Value: []byte{0x00, 0xff}},
		{Expected: math.MaxUint16, Value: []byte{0xff, 0xff}},
	}

	var buff bytes.Buffer
	_ = io.Writer(&buff)

	for _, test := range tests {

		buff.Write(test.Value)

		actual, err := ReadUnsignedShort(&buff)

		if err != nil {
			t.Error(err)
		}

		if actual != test.Expected {
			// Not equal
			t.Errorf("Unable to convert %v: %d != %d", test.Value, actual, test.Expected)
		}

		buff.Reset()
	}
}


func TestWriteVarInt(t *testing.T) {
	tests := []struct {
		Value    VarInt
		Expected []byte
	}{
		{Value: 0, Expected: []byte{0x00}},
		{Value: 1, Expected: []byte{0x01}},
		{Value: 2, Expected: []byte{0x02}},
		{Value: 127, Expected: []byte{0x7f}},
		{Value: 128, Expected: []byte{0x80, 0x01}},
		{Value: 255, Expected: []byte{0xff, 0x01}},
		{Value: 2147483647, Expected: []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
		{Value: -1, Expected: []byte{0xff, 0xff, 0xff, 0xff, 0x0f}},
		{Value: -2147483648, Expected: []byte{0x80, 0x80, 0x80, 0x80, 0x08}},
	}

	var buff bytes.Buffer
	_ = io.Writer(&buff)

	for _, test := range tests {
		err := WriteVarInt(&buff, test.Value)

		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(test.Expected, buff.Bytes()) != 0 {
			// Not equal
			t.Errorf("Unable to convert %d: %v != %v", test.Value, buff.Bytes(), test.Expected)
		}

		buff.Reset()
	}
}

func TestReadVarInt(t *testing.T) {
	tests := []struct {
		Expected VarInt
		Value    []byte
	}{
		{Expected: 0, Value: []byte{0x00}},
		{Expected: 1, Value: []byte{0x01}},
		{Expected: 2, Value: []byte{0x02}},
		{Expected: 127, Value: []byte{0x7f}},
		{Expected: 128, Value: []byte{0x80, 0x01}},
		{Expected: 255, Value: []byte{0xff, 0x01}},
		{Expected: 2147483647, Value: []byte{0xff, 0xff, 0xff, 0xff, 0x07}},
		{Expected: -1, Value: []byte{0xff, 0xff, 0xff, 0xff, 0x0f}},
		{Expected: -2147483648, Value: []byte{0x80, 0x80, 0x80, 0x80, 0x08}},
	}

	var buff bytes.Buffer
	_ = io.Writer(&buff)

	for _, test := range tests {
		buff.Write(test.Value)

		actual, err := ReadVarInt(&buff)

		if err != nil {
			t.Error(err)
		}

		if actual != test.Expected {
			// Not equal
			t.Errorf("Unable to convert %v: %d != %d", test.Value, actual, test.Expected)
		}

		buff.Reset()
	}
}

func TestWriteString(t *testing.T) {
	tests := []struct {
		Value    string
		Expected []byte
	}{
		{Value: "john", Expected: []byte{0x04, 0x6a, 0x6f, 0x68, 0x6e}},
		{Value: " doe ", Expected: []byte{0x05, 0x20, 0x64, 0x6f, 0x65, 0x20}},
		{Value: "😂😂😂", Expected: []byte{0x0c, 0xf0, 0x9f, 0x98, 0x82, 0xf0, 0x9f, 0x98, 0x82, 0xf0, 0x9f, 0x98, 0x82}},
		{Value: "(╯°Д°）╯︵/(.□ . \\)", Expected: []byte{0x1e, 0x28, 0xe2, 0x95, 0xaf, 0xc2, 0xb0, 0xd0, 0x94, 0xc2, 0xb0, 0xef, 0xbc, 0x89, 0xe2, 0x95, 0xaf, 0xef, 0xb8, 0xb5, 0x2f, 0x28, 0x2e, 0xe2, 0x96, 0xa1, 0x20, 0x2e, 0x20, 0x5c, 0x29}},
	}

	var buff bytes.Buffer
	_ = io.Writer(&buff)

	for _, test := range tests {
		err := WriteString(&buff, String(test.Value))

		if err != nil {
			t.Fatal(err)
		}

		if bytes.Compare(test.Expected, buff.Bytes()) != 0 {
			// Not equal
			t.Errorf(`Unable to convert "%s": %v != %v`, test.Value, buff.Bytes(), test.Expected)
		}

		buff.Reset()
	}
}

func TestReadString(t *testing.T) {
	tests := []struct {
		Expected String
		Value    []byte
	}{
		{Expected: "john", Value: []byte{0x04, 0x6a, 0x6f, 0x68, 0x6e}},
		{Expected: " doe ", Value: []byte{0x05, 0x20, 0x64, 0x6f, 0x65, 0x20}},
		{Expected: "😂😂😂", Value: []byte{0x0c, 0xf0, 0x9f, 0x98, 0x82, 0xf0, 0x9f, 0x98, 0x82, 0xf0, 0x9f, 0x98, 0x82}},
		{Expected: "(╯°Д°）╯︵/(.□ . \\)", Value: []byte{0x1e, 0x28, 0xe2, 0x95, 0xaf, 0xc2, 0xb0, 0xd0, 0x94, 0xc2, 0xb0, 0xef, 0xbc, 0x89, 0xe2, 0x95, 0xaf, 0xef, 0xb8, 0xb5, 0x2f, 0x28, 0x2e, 0xe2, 0x96, 0xa1, 0x20, 0x2e, 0x20, 0x5c, 0x29}},
	}

	var buff bytes.Buffer
	_ = io.Writer(&buff)

	for _, test := range tests {
		buff.Write(test.Value)

		actual, err := ReadString(&buff)

		if err != nil {
			t.Error(err)
		}

		if actual != test.Expected {
			// Not equal
			t.Errorf(`Unable to convert %v: "%s" != "%s"`, test.Value, actual, test.Expected)
		}

		buff.Reset()
	}
}
