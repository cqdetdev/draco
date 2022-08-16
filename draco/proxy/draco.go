package proxy

import (
	"errors"
	"sync"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sirupsen/logrus"
)

// Draco represents an instance of the Draco proxy.
type Draco struct {
	players chan *player.Player
	log     *logrus.Logger
	addr    string
}

// New returns a new Draco instance.
// If your server is using Dragonfly, be sure to use the Listener function instead.
func New(log *logrus.Logger, localAddr string) *Draco {
	return &Draco{
		players: make(chan *player.Player),
		addr:    localAddr,
		log:     log,
	}
}

// Start will start Draco! remoteAddr is the address of the target server, and localAddr is the address that players will connect to.
// Addresses should be formatted in the following format: "ip:port" (ex: "127.0.0.1:19132").
// If you're using dragonfly, use Listen instead of Start.
func (o *Draco) Start(remoteAddr string, resourcePackPath string, requirePacks bool) error {
	p, err := minecraft.NewForeignStatusProvider(remoteAddr)
	if err != nil {
		panic(err)
	}
	l, err := minecraft.ListenConfig{
		StatusProvider: p,
		// ResourcePacks:        utils.ResourcePacks(resourcePackPath),
		TexturePacksRequired: requirePacks,
	}.Listen("raknet", o.addr)
	if err != nil {
		return err
	}
	defer l.Close()
	o.log.Printf("Draco is now listening on %v and directing connections to %v!\n", o.addr, remoteAddr)
	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go o.handleConn(c.(*minecraft.Conn), l, remoteAddr)
	}
}

// handleConn handles a new incoming minecraft.Conn from the minecraft.Listener passed.
func (o *Draco) handleConn(conn *minecraft.Conn, listener *minecraft.Listener, remoteAddr string) {
	serverConn, err := minecraft.Dialer{
		IdentityData: conn.IdentityData(),
		ClientData:   conn.ClientData(),
	}.Dial("raknet", remoteAddr)
	if err != nil {
		return
	}

	data := serverConn.GameData()
	data.PlayerMovementSettings.MovementType = protocol.PlayerMovementModeServerWithRewind
	data.PlayerMovementSettings.RewindHistorySize = 60

	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(data); err != nil {
			return
		}
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			return
		}
		g.Done()
	}()
	g.Wait()

	g.Add(2)
	go func() {
		defer func() {
			_ = listener.Disconnect(conn, "connection lost")
			_ = serverConn.Close()
			g.Done()
		}()
		for {
			pk, err := conn.ReadPacket()
			if err != nil || conn == nil {
				return
			}

			if err := serverConn.WritePacket(pk); err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(conn, disconnect.Error())
				}
				return
			}
		}
	}()
	go func() {
		defer func() {
			_ = serverConn.Close()
			_ = listener.Disconnect(conn, "connection lost")
			g.Done()
		}()
		for {
			pk, err := serverConn.ReadPacket()
			if err != nil {
				if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
					_ = listener.Disconnect(conn, disconnect.Error())
				}
				return
			}

			if err := conn.WritePacket(pk); err != nil {
				return
			}
		}
	}()
	g.Wait()
	conn.Close()
}
