package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type PlayerAuthInputTranslator struct{}

func (PlayerAuthInputTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.PlayerAuthInput)
	earlier := &v486.PlayerAuthInput{
		Pitch: latest.Pitch,
		Yaw: latest.Yaw,
		Position: latest.Position,
		MoveVector: latest.MoveVector,
		HeadYaw: latest.HeadYaw,
		InputData: latest.InputData,
		InputMode: latest.InputMode,
		PlayMode: latest.PlayMode,
		GazeDirection: latest.GazeDirection,
		Tick: latest.Tick,
		Delta: latest.Delta,
		ItemInteractionData: latest.ItemInteractionData,
		ItemStackRequest: latest.ItemStackRequest,
		BlockActions: latest.BlockActions,
	}
	earlier.ItemInteractionData.HeldItem.Stack = translator.UpgradeItemStack(latest.ItemInteractionData.HeldItem.Stack)
	return earlier
}