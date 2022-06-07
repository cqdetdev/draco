package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryContentTranslator struct{}

func (InventoryContentTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.InventoryContent)
	items := make([]protocol.ItemInstance, 0, len(latest.Content))
	for _, it := range latest.Content {
		it.Stack = translator.DowngradeItemStack(it.Stack)
		items = append(items, it)
	}
	latest.Content = items
	return latest
}