package chunk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cqdetdev/draco/draco/latestmappings"
	"github.com/cqdetdev/draco/draco/state"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type (
	// Encoding is an encoding type used for Chunk encoding. Implementations of this interface are DiskEncoding and
	// NetworkEncoding, which can be used to encode a Chunk to an intermediate disk or network representation respectively.
	Encoding interface {
		encodePalette(buf *bytes.Buffer, p *Palette, e paletteEncoding)
		decodePalette(buf *bytes.Buffer, blockSize paletteSize, e paletteEncoding) (*Palette, error)
		network() byte
	}
	// paletteEncoding is an encoding type used for Chunk encoding. It is used to encode different types of palettes
	// (for example, blocks or biomes) differently.
	paletteEncoding interface {
		encode(buf *bytes.Buffer, v uint32)
		decode(buf *bytes.Buffer) (uint32, error)
	}
)

var (
	// NetworkEncoding is the Encoding used for sending a Chunk over network. It does not use NBT and writes varints.
	NetworkEncoding networkEncoding
	// BiomePaletteEncoding is the paletteEncoding used for encoding a palette of biomes.
	BiomePaletteEncoding biomePaletteEncoding
	// BlockPaletteEncoding is the paletteEncoding used for encoding a palette of block states encoded as NBT.
	BlockPaletteEncoding blockPaletteEncoding
)

// biomePaletteEncoding implements the encoding of biome palettes to disk.
type biomePaletteEncoding struct{}

func (biomePaletteEncoding) encode(buf *bytes.Buffer, v uint32) {
	_ = binary.Write(buf, binary.LittleEndian, v)
}
func (biomePaletteEncoding) decode(buf *bytes.Buffer) (uint32, error) {
	var v uint32
	return v, binary.Read(buf, binary.LittleEndian, &v)
}

// blockPaletteEncoding implements the encoding of block palettes to disk.
type blockPaletteEncoding struct{}

func (blockPaletteEncoding) encode(buf *bytes.Buffer, v uint32) {
	// Get the block state registered with the runtime IDs we have in the palette of the block storage
	// as we need the name and data value to store.
	name, props, _ := latestmappings.RuntimeIDToState(v)
	_ = nbt.NewEncoderWithEncoding(buf, nbt.LittleEndian).Encode(state.Block{Name: name, Properties: props, Version: CurrentBlockVersion})
}
func (blockPaletteEncoding) decode(buf *bytes.Buffer) (uint32, error) {
	var e state.Block
	if err := nbt.NewDecoderWithEncoding(buf, nbt.LittleEndian).Decode(&e); err != nil {
		return 0, fmt.Errorf("error decoding block palette entry: %w", err)
	}
	v, ok := latestmappings.StateToRuntimeID(e.Name, e.Properties)
	if !ok {
		return 0, fmt.Errorf("cannot get runtime ID of block state %v{%+v}", e.Name, e.Properties)
	}
	return v, nil
}

// networkEncoding implements the Chunk encoding for sending over network.
type networkEncoding struct{}

func (networkEncoding) network() byte { return 1 }
func (networkEncoding) encodePalette(buf *bytes.Buffer, p *Palette, _ paletteEncoding) {
	if p.size != 0 {
		_ = protocol.WriteVarint32(buf, int32(p.Len()))
	}
	for _, val := range p.values {
		_ = protocol.WriteVarint32(buf, int32(val))
	}
}
func (networkEncoding) decodePalette(buf *bytes.Buffer, blockSize paletteSize, _ paletteEncoding) (*Palette, error) {
	var paletteCount int32 = 1
	if blockSize > 0 {
		if err := protocol.Varint32(buf, &paletteCount); err != nil {
			return nil, fmt.Errorf("error reading palette entry count: %w", err)
		}
		if paletteCount <= 0 {
			return nil, fmt.Errorf("invalid palette entry count %v", paletteCount)
		}
	}

	blocks, temp := make([]uint32, paletteCount), int32(0)
	for i := int32(0); i < paletteCount; i++ {
		if err := protocol.Varint32(buf, &temp); err != nil {
			return nil, fmt.Errorf("error decoding palette entry: %w", err)
		}
		blocks[i] = uint32(temp)
	}
	return &Palette{values: blocks, size: blockSize}, nil
}
