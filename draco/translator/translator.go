package translator

import (
	"fmt"

	"github.com/cqdetdev/draco/draco/chunk"
	"github.com/cqdetdev/draco/draco/latestmappings"
	"github.com/cqdetdev/draco/draco/legacymappings"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Translator struct {
	Inbound map[uint32]TranslationHandler
	Outbound map[uint32]TranslationHandler
	Protocol int32
}

type TranslationHandler interface {
	Translate(p packet.Packet) packet.Packet
}

// dataKeyVariant is used for falling blocks and fake texts. This is necessary for falling block runtime ID translation.
const dataKeyVariant = 2

// downgradeSubChunk translates a 1.18.30 sub-chunk to a 1.18.12 one, updating all palette entries with the appropriate
// runtime IDs.
func DowngradeSubChunk(s *chunk.SubChunk) {
	for _, l := range s.Layers() {
		l.Palette().Replace(DowngradeBlockRuntimeID)
	}
}

// downgradeBlockRuntimeID translates a 1.18.30 runtime ID to a 1.18.12 one.
func DowngradeBlockRuntimeID(latestRID uint32) uint32 {
	name, properties, found := latestmappings.RuntimeIDToState(latestRID)
	if !found {
		panic(fmt.Errorf("downgrade block runtime id: could not find name for runtime id: %v", latestRID))
	}
	earlierRuntimeID, found := legacymappings.StateToRuntimeID(name, properties)
	if !found {
		panic(fmt.Errorf("downgrade block runtime id: could not find runtime id for name: %v", name))
	}
	return earlierRuntimeID
}

// upgradeBlockRuntimeID translates a 1.18.12 block runtime ID to a 1.18.30 one.
func UpgradeBlockRuntimeID(id uint32) uint32 {
	name, properties, found := legacymappings.RuntimeIDToState(id)
	if !found {
		panic(fmt.Errorf("upgrade block runtime id: could not find name for runtime id: %v", id))
	}
	latestRuntimeID, found := latestmappings.StateToRuntimeID(name, properties)
	if !found {
		panic(fmt.Errorf("upgrade block runtime id: could not find runtime id for name: %v", name))
	}
	return latestRuntimeID
}

// DowngradeEntityMetadata translates a 1.18.30 entity metadata to a 1.18.12 one.
func DowngradeEntityMetadata(metadata map[uint32]any) {
	if latestRID, ok := metadata[dataKeyVariant]; ok {
		metadata[dataKeyVariant] = int32(DowngradeBlockRuntimeID(uint32(latestRID.(int32))))
	}
}

// DowngradeRecipe downgrades a 1.18.30 recipe to a 1.18.12 one.
func DowngradeRecipe(latestInput []protocol.RecipeIngredientItem, latestOutput []protocol.ItemStack) ([]protocol.RecipeIngredientItem, []protocol.ItemStack) {
	input := make([]protocol.RecipeIngredientItem, 0, len(latestInput))
	output := make([]protocol.ItemStack, 0, len(latestOutput))
	for _, i := range latestInput {
		if i.Count > 0 {
			i.NetworkID = DowngradeItemRuntimeID(i.NetworkID)
		}
		input = append(input, i)
	}
	for _, o := range latestOutput {
		output = append(output, DowngradeItemStack(o))
	}
	return input, output
}

// DowngradeItemStack translates a 1.18.30 item stack to a 1.18.12 one, updating all palette entries with the appropriate
// runtime IDs.
func DowngradeItemStack(st protocol.ItemStack) protocol.ItemStack {
	if st.BlockRuntimeID > 0 {
		st.BlockRuntimeID = int32(DowngradeBlockRuntimeID(uint32(st.BlockRuntimeID)))
	}
	if st.HasNetworkID {
		st.NetworkID = DowngradeItemRuntimeID(st.NetworkID)
	}
	return st
}

// UpgradeItemStack translates a 1.18.12 item stack to a 1.18.30 one, updating all palette entries with the appropriate
// runtime IDs.
func UpgradeItemStack(st protocol.ItemStack) protocol.ItemStack {
	if st.BlockRuntimeID > 0 {
		st.BlockRuntimeID = int32(UpgradeBlockRuntimeID(uint32(st.BlockRuntimeID)))
	}
	if st.HasNetworkID {
		st.NetworkID = UpgradeItemRuntimeID(st.NetworkID)
	}
	return st
}

// UpgradeItemRuntimeID translates a 1.18.12 item runtime ID to a 1.18.30 one.
func UpgradeItemRuntimeID(latestRID int32) int32 {
	name, found := legacymappings.ItemRuntimeIDToName(latestRID)
	if !found {
		panic(fmt.Errorf("upgrade item runtime id: could not find name for runtime id: %v", latestRID))
	}
	earlierRuntimeID, found := latestmappings.ItemNameToRuntimeID(name)
	if !found {
		panic(fmt.Errorf("upgrade item runtime id: could not find runtime id for name: %v", name))
	}
	return earlierRuntimeID
}

// DowngradeItemRuntimeID translates a 1.18.30 item runtime ID to a 1.18.12 one.
func DowngradeItemRuntimeID(latestRID int32) int32 {
	name, found := latestmappings.ItemRuntimeIDToName(latestRID)
	if !found {
		panic(fmt.Errorf("downgrade item runtime id: could not find name for runtime id: %v", latestRID))
	}
	earlierRuntimeID, found := legacymappings.ItemNameToRuntimeID(name)
	if !found {
		panic(fmt.Errorf("downgrade item runtime id: could not find runtime id for name: %v", name))
	}
	return earlierRuntimeID
}
