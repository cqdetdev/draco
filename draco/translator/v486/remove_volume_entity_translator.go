package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type RemoveVolumeEntityTranslator struct{}

func (RemoveVolumeEntityTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.AddVolumeEntity)
	return &v486.RemoveVolumeEntity{EntityRuntimeID: latest.EntityRuntimeID}
}