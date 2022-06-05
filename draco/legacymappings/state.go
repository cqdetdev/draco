package legacymappings

import (
	"bytes"
	_ "embed"
	"github.com/cqdetdev/draco/draco/state"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte
	//go:embed block_aliases.nbt
	blockAliasesData []byte
	// stateRuntimeIDs holds a map for looking up the runtime ID of a block by the stateHash it produces.
	stateRuntimeIDs = map[state.Hash]uint32{}
	// runtimeIDToState holds a map for looking up the blockState of a block by its runtime ID.
	runtimeIDToState = map[uint32]state.Block{}
	// aliasMappings maps from a legacy block name alias to an updated name.
	aliasMappings = map[string]string{}
)

var (
	//go:embed item_runtime_ids.nbt
	itemRuntimeIDData []byte
	// itemRuntimeIDsToNames holds a map to translate item runtime IDs to string IDs.
	itemRuntimeIDsToNames = map[int32]string{}
	// itemNamesToRuntimeIDs holds a map to translate item string IDs to runtime IDs.
	itemNamesToRuntimeIDs = map[string]int32{}
)

// init initializes the item and state mappings.
func init() {
	var items map[string]int32
	if err := nbt.Unmarshal(itemRuntimeIDData, &items); err != nil {
		panic(err)
	}
	for name, rid := range items {
		itemNamesToRuntimeIDs[name] = rid
		itemRuntimeIDsToNames[rid] = name
	}

	var aliases map[string]string
	if err := nbt.Unmarshal(blockAliasesData, &aliases); err != nil {
		panic(err)
	}
	for alias, name := range aliases {
		aliasMappings[name] = alias
	}

	dec := nbt.NewDecoder(bytes.NewBuffer(blockStateData))

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	var s state.Block
	for {
		if err := dec.Decode(&s); err != nil {
			break
		}
		rid := uint32(len(stateRuntimeIDs))
		stateRuntimeIDs[state.HashBlock(s)] = rid
		runtimeIDToState[rid] = s
	}
}

// StateToRuntimeID converts a name and its state properties to a runtime ID.
func StateToRuntimeID(name string, properties map[string]any) (runtimeID uint32, found bool) {
	if updatedName, ok := aliasMappings[name]; ok {
		name = updatedName
	}
	rid, ok := stateRuntimeIDs[state.HashBlock(state.Block{Name: name, Properties: properties})]
	return rid, ok
}

// RuntimeIDToState converts a runtime ID to a name and its state properties.
func RuntimeIDToState(runtimeID uint32) (name string, properties map[string]any, found bool) {
	s := runtimeIDToState[runtimeID]
	return s.Name, s.Properties, true
}

// ItemRuntimeIDToName converts an item runtime ID to a string ID.
func ItemRuntimeIDToName(runtimeID int32) (name string, found bool) {
	name, ok := itemRuntimeIDsToNames[runtimeID]
	return name, ok
}

// ItemNameToRuntimeID converts a string ID to an item runtime ID.
func ItemNameToRuntimeID(name string) (runtimeID int32, found bool) {
	if updatedName, ok := aliasMappings[name]; ok {
		name = updatedName
	}
	rid, ok := itemNamesToRuntimeIDs[name]
	return rid, ok
}
