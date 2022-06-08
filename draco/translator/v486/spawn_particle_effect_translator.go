package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SpawnParticleEffectTranslator struct{}

func (SpawnParticleEffectTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.SpawnParticleEffect)
	return &v486.SpawnParticleEffect{
		Dimension:      latest.Dimension,
		EntityUniqueID: latest.EntityUniqueID,
		Position:       latest.Position,
		ParticleName:   latest.ParticleName,
	}
}