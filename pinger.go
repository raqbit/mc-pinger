package mcpinger

import (
	"context"
	"fmt"
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
func (p Pinger) Ping(address string) (*mcproto.ServerInfo, error) {
	return p.PingContext(context.Background(), address)
}

// PingTimeout acts like Ping but takes a timeout.
//
// If the timeout expires before the server info has been retrieved, an error is returned.
func (p Pinger) PingTimeout(address string, timeout time.Duration) (*mcproto.ServerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return p.PingContext(ctx, address)
}

// PingContext acts like Ping but takes a context.
//
// The provided Context must be non-nil. If the context expires before
// the server info has been retrieved, an error is returned.
func (p Pinger) PingContext(ctx context.Context, address string) (*mcproto.ServerInfo, error) {
	if ctx == nil {
		panic("nil context")
	}

	// Make TCP connection
	var d net.Dialer
	tcpConn, err := d.DialContext(ctx, "tcp", address)

	if err != nil {
		return nil, fmt.Errorf("could not connect to Minecraft server: %w", err)
	}

	// TODO: use ctx for all read/write actions

	// Wrap TCP connection in mcproto
	conn := mcproto.NewConnection(tcpConn, mcproto.ClientSide)
	defer conn.Close()

	// Send handshake
	if err = sendHandshakePacket(conn, address); err != nil {
		return nil, err
	}

	// Switch to status state
	conn.State = mcproto.StatusState

	// Send server query
	if err = sendServerQueryPacket(conn); err != nil {
		return nil, err
	}

	// Read server info response
	res, err := readServerInfoPacket(conn)

	if err != nil {
		return nil, err
	}

	return &res.Response, err
}

// Ping sends a server query to the specified address.
func Ping(address string) (*mcproto.ServerInfo, error) {
	var p Pinger
	return p.Ping(address)
}

// PingContext acts like Ping but takes a context.
//
// The provided Context must be non-nil. If the context expires before
// the server info has been retrieved, an error is returned.
func PingContext(ctx context.Context, address string) (*mcproto.ServerInfo, error) {
	var p Pinger
	return p.PingContext(ctx, address)
}

// PingTimeout acts like Ping but takes a timeout.
//
// If the timeout expires before the server info has been retrieved, an error is returned.
func PingTimeout(address string, timeout time.Duration) (*mcproto.ServerInfo, error) {
	var p Pinger
	return p.PingTimeout(address, timeout)
}

func sendHandshakePacket(conn *mcproto.Connection, address string) error {
	host, port, err := net.SplitHostPort(address)

	if err != nil {
		return fmt.Errorf("could not split host & port: %w", err)
	}

	serverPort, err := strconv.Atoi(port)

	if err != nil {
		return fmt.Errorf("could not parse port: %w", err)
	}

	handshakePkt := &mcproto.CHandshakePacket{
		ProtoVer:   -1,
		ServerAddr: host,
		ServerPort: uint16(serverPort),
		NextState:  mcproto.StatusState,
	}

	if err = conn.WritePacket(handshakePkt); err != nil {
		return fmt.Errorf("could not send handshake packet: %w", err)
	}

	return nil
}

func sendServerQueryPacket(w *mcproto.Connection) error {
	if err := w.WritePacket(&mcproto.CServerQueryPacket{}); err != nil {
		return fmt.Errorf("could not send server query request packet: %w", err)
	}

	return nil
}

func readServerInfoPacket(conn *mcproto.Connection) (*mcproto.SServerInfoPacket, error) {
	pkt, err := conn.ReadPacket()

	if err != nil {
		return nil, fmt.Errorf("could not read server info packet: %w", err)
	}

	resp, ok := pkt.(*mcproto.SServerInfoPacket)

	if !ok {
		return nil, fmt.Errorf("returned packet was not expected server info")
	}

	return resp, nil
}
