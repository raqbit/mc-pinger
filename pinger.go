package mcpinger

import (
	"bufio"
	"errors"
	"fmt"
	enc "github.com/Raqbit/mc-pinger/encoding"
	"github.com/Raqbit/mc-pinger/packet"
	"net"
	"strconv"
)

const (
	UnknownProtoVersion = -1
	StatusState         = 1
)

// Pinger allows you to retrieve server info.
type Pinger interface {
	Ping() (*ServerInfo, error)
}

type mcPinger struct {
	Host string
	Port uint16
}

// Error returned when the received packet type
// does not match the expected packet type.
type InvalidPacketError struct {
	expected enc.VarInt
	actual   enc.VarInt
}

func (i InvalidPacketError) Error() string {
	return fmt.Sprintf("Received invalid packet. Expected #%d, got #%d", i.expected, i.actual)
}

// Will connect to the Minecraft server,
// retrieve server status and return the server info.
func (p *mcPinger) Ping() (*ServerInfo, error) {

	address := net.JoinHostPort(p.Host, strconv.Itoa(int(p.Port)))

	// TODO: TIMEOUTS
	conn, err := net.Dial("tcp", address)

	rd := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	if err != nil {
		return nil, errors.New("could not connect to Minecraft server: " + err.Error())
	}

	defer conn.Close()

	err = p.sendHandshakePacket(w)

	if err != nil {
		return nil, err
	}

	err = p.sendRequestPacket(w)

	if err != nil {
		return nil, err
	}

	err = w.Flush()

	if err != nil {
		return nil, err
	}

	res, err := p.readPacket(rd)

	if err != nil {
		return nil, err
	}

	info, err := parseServerInfo([]byte(res.Json))

	return info, err
}

func (p *mcPinger) sendHandshakePacket(w *bufio.Writer) error {
	handshakePkt := &packet.HandshakePacket{
		ProtoVer:   UnknownProtoVersion,
		ServerAddr: enc.String(p.Host),
		ServerPort: enc.UnsignedShort(p.Port),
		NextState:  StatusState,
	}

	err := packet.WritePacket(handshakePkt, w)

	if err != nil {
		return errors.New("could not pack: " + err.Error())
	}

	return nil
}

func (p *mcPinger) sendRequestPacket(w *bufio.Writer) error {
	requestPkt := &packet.RequestPacket{}

	err := packet.WritePacket(requestPkt, w)

	if err != nil {
		return errors.New("could not pack: " + err.Error())
	}

	return nil
}

func (p *mcPinger) readPacket(rd *bufio.Reader) (*packet.ResponsePacket, error) {

	rp := &packet.ResponsePacket{}

	_, packetId, err := packet.ReadPacketHeader(rd)

	if packetId != rp.ID() {
		return nil, InvalidPacketError{expected: rp.ID(), actual: packetId}
	}

	if err != nil {
		return nil, err
	}

	err = rp.Unmarshal(rd)

	if err != nil {
		return nil, err
	}

	return rp, nil
}

// Creates a new Pinger with specified host & port
// to connect to a minecraft server
func New(host string, port uint16) Pinger {
	return &mcPinger{
		Host: host,
		Port: port,
	}
}
