package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var Translator translator.Translator = translator.Translator{
	Inbound: map[uint32]translator.TranslationHandler{
	},
	Outbound: map[uint32]translator.TranslationHandler{
		packet.IDMobEquipment: MobEquipmentTranslator{},
	},
	Protocol: 486,
}