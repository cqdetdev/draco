package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventorySlotTranslator struct{}

func (InventorySlotTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.InventorySlot)
	latest.NewItem.Stack = translator.DowngradeItemStack(latest.NewItem.Stack)
	return latest
}