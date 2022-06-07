package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type UpdateBlockTranslator struct{}

func (UpdateBlockTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.UpdateBlock)
	latest.NewBlockRuntimeID = translator.DowngradeBlockRuntimeID(latest.NewBlockRuntimeID)
	return latest
}