package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddActorTranslator struct{}

func (AddActorTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.AddActor)
	translator.DowngradeEntityMetadata(latest.EntityMetadata)
	return latest
}