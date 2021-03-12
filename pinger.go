package mcpinger

import (
	"context"
	"fmt"
	"github.com/Raqbit/mcproto/encoding"
	"github.com/Raqbit/mcproto/packet"
	"github.com/Raqbit/mcproto/types"
	"net"
	"strconv"
	"time"

	"github.com/Raqbit/mcproto"
)

// A Pinger contains options for pinging a Minecraft server.
//
// The zero value for each field is equivalent to pinging
// without that option. Pinging with the zero value of Pinger
// is therefore equivalent to just calling the Ping function.
//
// It is safe to call Pinger's methods concurrently.
type Pinger struct {
}

// Ping sends a server query to the specified address.
func (p Pinger) Ping(host, port string) (*packet.ServerInfo, error) {
	return p.PingContext(context.Background(), host, port)
}

// PingTimeout acts like Ping but takes a timeout.
//
// If the timeout expires before the server info has been retrieved, an error is returned.
func (p Pinger) PingTimeout(host, port string, timeout time.Duration) (*packet.ServerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return p.PingContext(ctx, host, port)
}

// PingContext acts like Ping but takes a context.
//
// The provided Context must be non-nil. If the context expires before
// the server info has been retrieved, an error is returned.
func (p Pinger) PingContext(ctx context.Context, host, port string) (*packet.ServerInfo, error) {
	if ctx == nil {
		panic("nil context")
	}

	// Connect to Minecraft server
	conn, address, err := mcproto.DialContext(ctx, host, port)

	if err != nil {
		return nil, fmt.Errorf("could not connect to Minecraft server: %w", err)
	}

	defer conn.Close()

	// Send handshake
	if err = sendHandshakePacket(ctx, conn, address); err != nil {
		return nil, err
	}

	// Switch to the state we requested
	conn.SetState(types.ConnectionStateStatus)

	// Send server query
	if err = sendServerQueryPacket(ctx, conn); err != nil {
		return nil, err
	}

	// Read server info response
	res, err := readServerInfoPacket(ctx, conn)

	if err != nil {
		return nil, err
	}

	return &res.Response, err
}

// Ping sends a server query to the specified address.
func Ping(host, port string) (*packet.ServerInfo, error) {
	var p Pinger
	return p.Ping(host, port)
}

// PingContext acts like Ping but takes a context.
//
// The provided Context must be non-nil. If the context expires before
// the server info has been retrieved, an error is returned.
func PingContext(ctx context.Context, host, port string) (*packet.ServerInfo, error) {
	var p Pinger
	return p.PingContext(ctx, host, port)
}

// PingTimeout acts like Ping but takes a timeout.
//
// If the timeout expires before the server info has been retrieved, an error is returned.
func PingTimeout(host, port string, timeout time.Duration) (*packet.ServerInfo, error) {
	var p Pinger
	return p.PingTimeout(host, port, timeout)
}

func sendHandshakePacket(ctx context.Context, conn mcproto.Connection, address string) error {
	host, port, err := net.SplitHostPort(address)

	if err != nil {
		return fmt.Errorf("could not split host & port: %w", err)
	}

	serverPort, err := strconv.Atoi(port)

	if err != nil {
		return fmt.Errorf("could not parse port: %w", err)
	}

	handshakePkt := &packet.HandshakePacket{
		ProtoVer:   -1,
		ServerAddr: encoding.String(host),
		ServerPort: encoding.UnsignedShort(uint16(serverPort)),
		NextState:  types.ConnectionStateStatus,
	}

	if err = conn.WritePacket(ctx, handshakePkt); err != nil {
		return fmt.Errorf("could not send handshake packet: %w", err)
	}

	return nil
}

func sendServerQueryPacket(ctx context.Context, w mcproto.Connection) error {
	if err := w.WritePacket(ctx, &packet.ServerQueryPacket{}); err != nil {
		return fmt.Errorf("could not send server query request packet: %w", err)
	}

	return nil
}

func readServerInfoPacket(ctx context.Context, conn mcproto.Connection) (*packet.ServerInfoPacket, error) {
	pkt, err := conn.ReadPacket(ctx)

	if err != nil {
		return nil, fmt.Errorf("could not read server info packet: %w", err)
	}

	resp, ok := pkt.(*packet.ServerInfoPacket)

	if !ok {
		return nil, fmt.Errorf("returned packet was not expected server info")
	}

	return resp, nil
}
