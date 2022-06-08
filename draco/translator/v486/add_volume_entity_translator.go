package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddVolumeEntityTranslator struct{}

func (AddVolumeEntityTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.AddVolumeEntity)
	return &v486.AddVolumeEntity{
		EntityRuntimeID:    latest.EntityRuntimeID,
		EntityMetadata:     latest.EntityMetadata,
		EncodingIdentifier: latest.EncodingIdentifier,
		InstanceIdentifier: latest.InstanceIdentifier,
		EngineVersion:      latest.EngineVersion,
	}
}