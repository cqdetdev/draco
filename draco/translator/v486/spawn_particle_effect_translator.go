package v486

import (
	"github.com/cqdetdev/draco/draco/legacy"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SpawnParticleEffectTranslator struct{}

func (SpawnParticleEffectTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.SpawnParticleEffect)
	return &legacy.SpawnParticleEffect{
		Dimension:      latest.Dimension,
		EntityUniqueID: latest.EntityUniqueID,
		Position:       latest.Position,
		ParticleName:   latest.ParticleName,
	}
}