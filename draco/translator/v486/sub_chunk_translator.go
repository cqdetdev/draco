package v486

import (
	"bytes"

	"github.com/cqdetdev/draco/draco/chunk"
	"github.com/cqdetdev/draco/draco/latestmappings"
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SubChunkTranslator struct{}

func (SubChunkTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.SubChunk)
	entries := make([]protocol.SubChunkEntry, 0, len(latest.SubChunkEntries))
	for _, e := range latest.SubChunkEntries {
		if e.Result == protocol.SubChunkResultSuccess {
			var ind uint8
			buf := bytes.NewBuffer(e.RawPayload)
			air, _ := latestmappings.StateToRuntimeID("minecraft:air", nil)
			s, err := chunk.DecodeSubChunk(air, cube.Range{-64, 319}, buf, &ind, chunk.NetworkEncoding)
			if err != nil {
				panic(err)
			}
			translator.DowngradeSubChunk(s)
			serialisedSubChunk := chunk.EncodeSubChunk(s, chunk.NetworkEncoding, cube.Range{-64, 319}, int(ind))
			e.RawPayload = append(serialisedSubChunk, buf.Bytes()...)
		}
		entries = append(entries, e)
	}
	latest.SubChunkEntries = entries
	return latest
}