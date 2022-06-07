package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SetActorDataTranslator struct{}

func (SetActorDataTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.SetActorData)
	translator.DowngradeEntityMetadata(latest.EntityMetadata)
	return latest
}