package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type MobEquipmentTranslator struct{}

func (MobEquipmentTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.MobEquipment)
	latest.NewItem.Stack = translator.UpgradeItemStack(latest.NewItem.Stack)
	return latest

}