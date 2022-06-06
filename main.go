package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	// "sync"

	"github.com/cqdetdev/draco/draco"
	"github.com/pelletier/go-toml"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"golang.org/x/oauth2"
)

// The following program implements a proxy that forwards players from one local address to a remote address.
func main() {
	l := log.Default()
	c := readConfig()
	if err := draco.InitializeToken(l); err != nil {
		log.Fatal(err)
	}

	p, err := minecraft.NewForeignStatusProvider(c.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}

	tempConn, err := minecraft.Dialer{
		TokenSource: draco.TokenSrc,
	}.Dial("raknet", c.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}
	_ = tempConn.Close()

	li, err := minecraft.ListenConfig{
		AcceptedProtocols: []minecraft.Protocol{
			draco.Protocol{},
		},
		StatusProvider: p,
		ResourcePacks:  tempConn.ResourcePacks(),
	}.Listen("raknet", c.Connection.LocalAddress)
	if err != nil {
		panic(err)
	}

	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			panic(err)
		}

		go handleConn(conn.(*minecraft.Conn), li, c, draco.TokenSrc)
	}
}

func handleConn(conn *minecraft.Conn, listener *minecraft.Listener, c config, src oauth2.TokenSource) {
	serverConn, err := minecraft.Dialer{
		TokenSource: src,
		ClientData:  conn.ClientData(),
		// TODO: Properly support the client cache.
	}.Dial("raknet", c.Connection.RemoteAddress)
	if err != nil {
		panic(err)
	}

	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if err := conn.StartGame(serverConn.GameData()); err != nil {
			panic(err)
		}
		g.Done()
	}()
	go func() {
		if err := serverConn.DoSpawn(); err != nil {
			panic(err)
		}
		g.Done()
	}()
	g.Wait()

	go func() {
		defer listener.Disconnect(conn, "connection lost")
		defer serverConn.Close()
		for {
			pk, err := conn.ReadPacket()
			if err != nil {
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
		defer serverConn.Close()
		defer listener.Disconnect(conn, "connection lost")
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
}

type config struct {
	Connection struct {
		LocalAddress  string
		RemoteAddress string
	}
}

func readConfig() config {
	c := config{}
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		f, err := os.Create("config.toml")
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
		data, err := toml.Marshal(c)
		if err != nil {
			log.Fatalf("error encoding default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default config: %v", err)
		}
		_ = f.Close()
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}
	if c.Connection.LocalAddress == "" {
		c.Connection.LocalAddress = "0.0.0.0:19132"
	}
	data, _ = toml.Marshal(c)
	if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
		log.Fatalf("error writing config file: %v", err)
	}
	return c
}

func packetHandle(header packet.Header, payload []byte, src net.Addr, dst net.Addr) {
	log := fmt.Sprintf("%v -> %v (0x%X)\n", src, dst, header.PacketID)
	f, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString(log)
	f.Close()
}
