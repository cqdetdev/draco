package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type CraftingDataTranslator struct{}

func (CraftingDataTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.CraftingData)
	recipes := make([]protocol.Recipe, 0, len(latest.Recipes))
	for _, r := range latest.Recipes {
		switch r := r.(type) {
		case *protocol.ShapedRecipe:
			r.Input, r.Output = translator.DowngradeRecipe(r.Input, r.Output)
		case *protocol.ShapelessRecipe:
			r.Input, r.Output = translator.DowngradeRecipe(r.Input, r.Output)
		}
		recipes = append(recipes, r)
	}
	latest.Recipes = recipes
	return latest
}