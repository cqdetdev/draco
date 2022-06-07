package v486

import (
	"github.com/cqdetdev/draco/draco/legacy"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddVolumeEntityTranslator struct{}

func (AddVolumeEntityTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.AddVolumeEntity)
	return &legacy.AddVolumeEntity{
		EntityRuntimeID:    latest.EntityRuntimeID,
		EntityMetadata:     latest.EntityMetadata,
		EncodingIdentifier: latest.EncodingIdentifier,
		InstanceIdentifier: latest.InstanceIdentifier,
		EngineVersion:      latest.EngineVersion,
	}
}