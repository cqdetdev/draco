package v486

import (
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type AddPlayerTranslator struct{}

func (AddPlayerTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.AddPlayer)
	earlier := &v486.AddPlayer{
		UUID:                    latest.UUID,
		Username:                latest.Username,
		EntityUniqueID:          latest.EntityUniqueID,
		EntityRuntimeID:         latest.EntityRuntimeID,
		PlatformChatID:          latest.PlatformChatID,
		Position:                latest.Position,
		Velocity:                latest.Velocity,
		Pitch:                   latest.Pitch,
		Yaw:                     latest.Yaw,
		HeadYaw:                 latest.HeadYaw,
		HeldItem:                latest.HeldItem,
		EntityMetadata:          latest.EntityMetadata,
		Flags:                   latest.Flags,
		CommandPermissionLevel:  latest.CommandPermissionLevel,
		ActionPermissions:       latest.ActionPermissions,
		PermissionLevel:         latest.PermissionLevel,
		CustomStoredPermissions: latest.CustomStoredPermissions,
		PlayerUniqueID:          latest.PlayerUniqueID,
		EntityLinks:             latest.EntityLinks,
		DeviceID:                latest.DeviceID,
		BuildPlatform:           latest.BuildPlatform,
	}
	earlier.HeldItem.Stack = translator.DowngradeItemStack(latest.HeldItem.Stack)
	return earlier
}