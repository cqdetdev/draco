package state

import (
	"fmt"
	"sort"
	"strings"
	"unsafe"
)

// Block holds a combination of a name and properties, together with a version.
type Block struct {
	// Name is the name of the block.
	Name string `nbt:"name"`
	// Properties is a map of properties that define the block's state.
	Properties map[string]interface{} `nbt:"states"`
	// Version is the version of the block state.
	Version int32 `nbt:"version"`
}

// Hash is a struct that may be used as a map key for block states. It contains the name of the block state
// and an encoded version of the properties.
type Hash struct {
	Name, Properties string
}

// HashBlock produces a Hash for the Block given.
func HashBlock(state Block) Hash {
	hash := Hash{Name: state.Name}
	if state.Properties == nil {
		// If the properties are nil, we don't need to hash them.
		return hash
	}

	keys := make([]string, 0, len(state.Properties))
	for k := range state.Properties {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	var b strings.Builder
	for _, k := range keys {
		switch v := state.Properties[k].(type) {
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

	hash.Properties = b.String()
	return hash
}
