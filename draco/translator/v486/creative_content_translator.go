package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type CreativeContentTranslator struct{}

func (CreativeContentTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.CreativeContent)
	items := make([]protocol.CreativeItem, 0, len(latest.Items))
	for _, it := range latest.Items {
		it.Item = translator.DowngradeItemStack(it.Item)
		items = append(items, it)
	}
	latest.Items = items
	return latest
}