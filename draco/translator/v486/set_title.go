package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type SetTitleTranslator struct{}

func (SetTitleTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.SetTitle)
	earlier := &v486.SetTitle{
		ActionType: latest.ActionType,
		Text: latest.Text,
		FadeInDuration: latest.FadeInDuration,
		RemainDuration: latest.RemainDuration,
		FadeOutDuration: latest.FadeOutDuration,
	}
	return earlier
}