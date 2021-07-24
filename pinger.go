package mcpinger

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/pires/go-proxyproto"
	"net"
	"strconv"
	"time"

	enc "github.com/Raqbit/mc-pinger/encoding"
	"github.com/Raqbit/mc-pinger/packet"
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
	Host    string
	Port    uint16
	Context context.Context
	Timeout time.Duration

	UseProxy     bool
	ProxyVersion byte
}

// InvalidPacketError returned when the received packet type
// does not match the expected packet type.
type InvalidPacketError struct {
	expected enc.VarInt
	actual   enc.VarInt
}

func (i InvalidPacketError) Error() string {
	return fmt.Sprintf("Received invalid packet. Expected #%d, got #%d", i.expected, i.actual)
}

func (p *mcPinger) Ping() (*ServerInfo, error) {
	if p.Timeout > 0 && p.Context == nil {
		ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
		p.Context = ctx
		defer cancel()
	}
	if p.Context == nil {
		p.Context = context.Background()
	}
	return p.ping()
}

// Will connect to the Minecraft server,
// retrieve server status and return the server info.
func (p *mcPinger) ping() (*ServerInfo, error) {

	if p.Context == nil {
		panic("Context is nil!")
	}

	address := net.JoinHostPort(p.Host, strconv.Itoa(int(p.Port)))

	var d net.Dialer

	conn, err := d.DialContext(p.Context, "tcp", address)

	if err != nil {
		return nil, errors.New("could not connect to Minecraft server: " + err.Error())
	}

	if p.UseProxy {
		err = p.writeProxyHeader(conn)
		if err != nil {
			return nil, errors.New("could not write PROXY header: " + err.Error())
		}
	}

	rd := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

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

	_, packetID, err := packet.ReadPacketHeader(rd)

	if packetID != rp.ID() {
		return nil, InvalidPacketError{expected: rp.ID(), actual: packetID}
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

func (p *mcPinger) writeProxyHeader(conn net.Conn) error {
	header := proxyproto.HeaderProxyFromAddrs(p.ProxyVersion, conn.LocalAddr(), conn.RemoteAddr())

	_, err := header.WriteTo(conn)
	return err
}

// New Creates a new Pinger with specified host & port
// to connect to a minecraft server
func New(host string, port uint16, options ...McPingerOption) Pinger {
	p := &mcPinger{
		Host: host,
		Port: port,
	}
	for _, opt := range options {
		opt(p)
	}
	return p
}

// NewTimed Creates a new Pinger with specified host & port
// to connect to a minecraft server with Timeout
func NewTimed(host string, port uint16, timeout time.Duration) Pinger {
	return &mcPinger{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
}

// NewContext Creates a new Pinger with the given Context
func NewContext(ctx context.Context, host string, port uint16) Pinger {
	return &mcPinger{
		Host:    host,
		Port:    port,
		Context: ctx,
	}
}

// McPingerOption instances can be combined when creating a new Pinger
type McPingerOption func(p *mcPinger)

func WithTimeout(timeout time.Duration) McPingerOption {
	return func(p *mcPinger) {
		p.Timeout = timeout
	}
}

func WithContext(ctx context.Context) McPingerOption {
	return func(p *mcPinger) {
		p.Context = ctx
	}
}

// WithProxyProto enables support for Bungeecord's proxy_protocol feature, which listens for
// PROXY protocol connections via HAproxy. version must be 1 (text) or 2 (binary).
func WithProxyProto(version byte) McPingerOption {
	return func(p *mcPinger) {
		p.UseProxy = true
		p.ProxyVersion = version
	}
}
