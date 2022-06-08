package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ResourcePacksInfoTranslator struct{}

func (ResourcePacksInfoTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.ResourcePacksInfo)
	earlier := &v486.ResourcePacksInfo{
		TexturePackRequired: latest.TexturePackRequired,
		HasScripts: latest.HasScripts,
		BehaviourPacks: latest.BehaviourPacks,
		TexturePacks: latest.TexturePacks,
	}
	return earlier
}