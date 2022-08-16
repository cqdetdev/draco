package player

import (
	"context"
	"net"
	"time"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Player is a wrapper struct for users
// that are proxied within Draco.
type Player struct {
	conn *minecraft.Conn
}

func NewPlayer(conn *minecraft.Conn) *Player {
	return &Player{
		conn: conn,
	}
}

func (p *Player) Conn() *minecraft.Conn {
	return p.conn
}

// IdentityData returns the login.IdentityData of a player. It contains the UUID, XUID and username of the connection.
func (p *Player) IdentityData() login.IdentityData {
	return p.conn.IdentityData()
}

// ClientData returns the login.ClientData of a player. This includes less sensitive data of the player like its skin,
// language code and other non-essential information.
func (p *Player) ClientData() login.ClientData {
	return p.conn.ClientData()
}

// ClientCacheEnabled specifies if the conn has the client cache, used for caching chunks client-side, enabled or
// not. Some platforms, like the Nintendo Switch, have this disabled at all times.
func (p *Player) ClientCacheEnabled() bool {
	// todo: support client cache
	//return p.conn.ClientCacheEnabled()
	return false
}

// ChunkRadius returns the chunk radius as requested by the client at the other end of the conn.
func (p *Player) ChunkRadius() int {
	return p.conn.ChunkRadius()
}

// Latency returns the current latency measured over the conn.
func (p *Player) Latency() time.Duration {
	return p.conn.Latency()
}

// Flush flushes the packets buffered by the conn, sending all of them out immediately.
func (p *Player) Flush() error {
	return p.conn.Flush()
}

// RemoteAddr returns the remote network address.
func (p *Player) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// WritePacket will call minecraft.Conn.WritePacket and process the packet with oomph.
func (p *Player) WritePacket(pk packet.Packet) error {
	if err := p.conn.WritePacket(pk); err != nil {
		return err
	}
	return nil
}

// ReadPacket will call minecraft.Conn.ReadPacket and process the packet with oomph.
func (p *Player) ReadPacket() (pk packet.Packet, err error) {
	if pk, err = p.conn.ReadPacket(); err != nil {
		return nil, err
	}
	return pk, err
}

// StartGameContext starts the game for the conn with a context to cancel it.
func (p *Player) StartGameContext(ctx context.Context, data minecraft.GameData) error {
	data.PlayerMovementSettings.MovementType = protocol.PlayerMovementModeServerWithRewind
	return p.conn.StartGameContext(ctx, data)
}

// Close closes the player.
func (p *Player) Close() error {
	_ = p.conn.Close()
	return nil
}
