package v486

import (
	v486 "github.com/cqdetdev/draco/draco/translator/v486"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Protocol is the protocol used to support the Minecraft 1.18.10 protocol (486).
type Protocol struct {
	minecraft.Protocol
}

// ID ...
func (Protocol) ID() int32 {
	return 486
}

// Ver ...
func (Protocol) Ver() string {
	return "1.18.10"
}

// Packets ...
func (p Protocol) Packets() packet.Pool {
	return packet.NewPool()
}

// ConvertToLatest ...
func (p Protocol) ConvertToLatest(pk packet.Packet) packet.Packet {
	if t, ok := v486.Translator.Outbound[pk.ID()]; ok {
		return t.Translate(pk)
	}
	return pk
}

// ConvertFromLatest ...
func (p Protocol) ConvertFromLatest(pk packet.Packet) packet.Packet {
	if t, ok := v486.Translator.Inbound[pk.ID()]; ok {
		return t.Translate(pk)
	}
	return pk
}