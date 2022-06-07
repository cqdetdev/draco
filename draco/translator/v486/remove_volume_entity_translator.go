package v486

import (
	"github.com/cqdetdev/draco/draco/legacy"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type RemoveVolumeEntityTranslator struct{}

func (RemoveVolumeEntityTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.AddVolumeEntity)
	return &legacy.RemoveVolumeEntity{EntityRuntimeID: latest.EntityRuntimeID}
}