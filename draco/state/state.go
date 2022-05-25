package state

import (
	"bytes"
	_ "embed"
	"fmt"
	"math"
	"sort"
	"strings"
	"unsafe"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity/physics"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed block_states.nbt
	blockStateData []byte
	// blocks holds a list of all registered Blocks indexed by their runtime ID. Blocks that were not explicitly
	// registered are of the type unknownBlock.
	blocks []world.Block
	// stateRuntimeIDs holds a map for looking up the runtime ID of a block by the stateHash it produces.
	stateRuntimeIDs = map[stateHash]uint32{}
	// nbtBlocks holds a list of NBTer implementations for blocks registered that implement the NBTer interface.
	// These are indexed by their runtime IDs. Blocks that do not implement NBTer have a false value in this slice.
	nbtBlocks []bool
	// randomTickBlocks holds a list of RandomTicker implementations for blocks registered that implement the RandomTicker interface.
	// These are indexed by their runtime IDs. Blocks that do not implement RandomTicker have a false value in this slice.
	randomTickBlocks []bool
	// airRID is the runtime ID of an air block.
	airRID uint32
)

// itemHash is a combination of an item's name and metadata. It is used as a key in hash maps.
type itemHash struct {
	name string
	meta int16
}

var (
	//go:embed item_runtime_ids.nbt
	itemRuntimeIDData []byte
	// items holds a list of all registered items, indexed using the itemHash created when calling
	// Item.EncodeItem.
	items = map[itemHash]world.Item{}
	// customItems holds a list of all registered custom items.
	customItems []world.CustomItem
	// itemRuntimeIDsToNames holds a map to translate item runtime IDs to string IDs.
	itemRuntimeIDsToNames = map[int32]string{}
	// itemNamesToRuntimeIDs holds a map to translate item string IDs to runtime IDs.
	itemNamesToRuntimeIDs = map[string]int32{}
)

var StateToRuntimeID func(string, map[string]interface{}) (uint32, bool)
var RuntimeIDToState func(runtimeID uint32) (name string, properties map[string]any, found bool)

func init() {
	var m map[string]int32
	err := nbt.Unmarshal(itemRuntimeIDData, &m)

	if err != nil {
		panic(err)
	}
	for name, rid := range m {
		itemNamesToRuntimeIDs[name] = rid
		itemRuntimeIDsToNames[rid] = name
	}

	dec := nbt.NewDecoder(bytes.NewBuffer(blockStateData))

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	var s blockState
	for {
		if err := dec.Decode(&s); err != nil {
			break
		}
		registerBlockState(s)
	}

	RuntimeIDToState = func(runtimeID uint32) (name string, properties map[string]any, found bool) {
		if runtimeID >= uint32(len(blocks)) {
			return "", nil, false
		}
		name, properties = blocks[runtimeID].EncodeBlock()
		return name, properties, true
	}

	StateToRuntimeID = func(name string, properties map[string]interface{}) (runtimeID uint32, found bool) {
		switch name {
		case "minecraft:invisible_bedrock":
			name = "minecraft:invisibleBedrock"
		case "minecraft:sea_lantern":
			name = "minecraft:seaLantern"
		case "minecraft:concrete_powder":
			name = "minecraft:concretePowder"
		}
		rid, ok := stateRuntimeIDs[stateHash{name: name, properties: hashProperties(properties)}]
		return rid, ok
	}
}

// registerBlockState registers a new blockState to the states slice. The function panics if the properties the
// blockState hold are invalid or if the blockState was already registered.
func registerBlockState(s blockState) {
	h := stateHash{name: s.Name, properties: hashProperties(s.Properties)}
	if _, ok := stateRuntimeIDs[h]; ok {
		panic(fmt.Sprintf("cannot register the same state twice (%+v)", s))
	}
	rid := uint32(len(blocks))
	if s.Name == "minecraft:air" {
		airRID = rid
	}
	stateRuntimeIDs[h] = rid
	blocks = append(blocks, unknownBlock{s})

	nbtBlocks = append(nbtBlocks, false)
	randomTickBlocks = append(randomTickBlocks, false)
	// chunk.FilteringBlocks = append(chunk.FilteringBlocks, 15)
	// chunk.LightBlocks = append(chunk.LightBlocks, 0)
}

// unknownModel is the model used for unknown blocks. It is the equivalent of a fully solid model.
type unknownModel struct{}

// AABB ...
func (u unknownModel) AABB(cube.Pos, *world.World) []physics.AABB {
	return []physics.AABB{physics.NewAABB(mgl64.Vec3{}, mgl64.Vec3{1, 1, 1})}
}

// FaceSolid ...
func (u unknownModel) FaceSolid(cube.Pos, cube.Face, *world.World) bool {
	return true
}

// unknownBlock represents a block that has not yet been implemented. It is used for registering block
// states that haven't yet been added.
type unknownBlock struct {
	blockState
}

// EncodeBlock ...
func (b unknownBlock) EncodeBlock() (string, map[string]interface{}) {
	return b.Name, b.Properties
}

// Model ...
func (unknownBlock) Model() world.BlockModel {
	return unknownModel{}
}

// Hash ...
func (b unknownBlock) Hash() uint64 {
	return math.MaxUint64
}

// blockState holds a combination of a name and properties, together with a version.
type blockState struct {
	Name       string                 `nbt:"name"`
	Properties map[string]interface{} `nbt:"states"`
	Version    int32                  `nbt:"version"`
}

// stateHash is a struct that may be used as a map key for block states. It contains the name of the block state
// and an encoded version of the properties.
type stateHash struct {
	name, properties string
}

// HashProperties produces a hash for the block properties held by the blockState.
func hashProperties(properties map[string]interface{}) string {
	if properties == nil {
		return ""
	}
	keys := make([]string, 0, len(properties))
	for k := range properties {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	var b strings.Builder
	for _, k := range keys {
		switch v := properties[k].(type) {
		case bool:
			if v {
				b.WriteByte(1)
			} else {
				b.WriteByte(0)
			}
		case uint8:
			b.WriteByte(v)
		case int32:
			a := *(*[4]byte)(unsafe.Pointer(&v))
			b.Write(a[:])
		case string:
			b.WriteString(v)
		default:
			// If block encoding is broken, we want to find out as soon as possible. This saves a lot of time
			// debugging in-game.
			panic(fmt.Sprintf("invalid block property type %T for property %v", v, k))
		}
	}

	return b.String()
}

// ItemByName attempts to return an item by a name and a metadata value.
func ItemByName(name string, meta int16) (world.Item, bool) {
	it, ok := items[itemHash{name: name, meta: meta}]
	if !ok {
		// Also try obtaining the item with a metadata value of 0, for cases with durability.
		it, ok = items[itemHash{name: name}]
	}
	return it, ok
}

// ItemRuntimeID attempts to return the runtime ID of the Item passed. False is returned if the Item is not
// registered.
func ItemRuntimeID(i world.Item) (rid int32, meta int16, ok bool) {
	name, meta := i.EncodeItem()
	rid, ok = itemNamesToRuntimeIDs[name]
	return rid, meta, ok
}

// ItemByRuntimeID attempts to return an Item by the runtime ID passed. If no item with that runtime ID exists,
// false is returned. ItemByRuntimeID also tries to find the item with a metadata value of 0.
func ItemByRuntimeID(rid int32, meta int16) (world.Item, bool) {
	name, ok := itemRuntimeIDsToNames[rid]
	if !ok {
		return nil, false
	}
	return ItemByName(name, meta)
}

// Items returns a slice of all registered items.
func Items() []world.Item {
	m := make([]world.Item, 0, len(items))
	for _, i := range items {
		m = append(m, i)
	}
	return m
}

// CustomItems returns a slice of all registered custom items.
func CustomItems() []world.CustomItem {
	return customItems
}
