package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerAuthInputTranslator struct{}

func (PlayerAuthInputTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.PlayerAuthInput)
	latest.ItemInteractionData.HeldItem.Stack = translator.UpgradeItemStack(latest.ItemInteractionData.HeldItem.Stack)
	return latest
}