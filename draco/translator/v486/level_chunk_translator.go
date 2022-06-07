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

type LevelChunkTranslator struct{}

func (LevelChunkTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.LevelChunk)
	if latest.SubChunkRequestMode == protocol.SubChunkRequestModeLegacy {
		readBuf := bytes.NewBuffer(latest.RawPayload)
		air, _ := latestmappings.StateToRuntimeID("minecraft:air", nil)
		c, err := chunk.NetworkDecode(air, readBuf, int(latest.SubChunkCount), cube.Range{-64, 319})
		if err != nil {
			panic(err)
		}
		for _, s := range c.Sub() {
			translator.DowngradeSubChunk(s)
		}

		writeBuf, data := bytes.NewBuffer(nil), chunk.Encode(c, chunk.NetworkEncoding)
		for i := range data.SubChunks {
			_, _ = writeBuf.Write(data.SubChunks[i])
		}
		_, _ = writeBuf.Write(data.Biomes)

		latest.RawPayload = append(writeBuf.Bytes(), readBuf.Bytes()...)
	}
	return latest
}