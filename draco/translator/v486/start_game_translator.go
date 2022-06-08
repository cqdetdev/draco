package v486

import (
	"github.com/cqdetdev/draco/draco/legacymappings"
	v486 "github.com/cqdetdev/draco/draco/packet/v486"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type StartGameTranslator struct{}

func (StartGameTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.StartGame)
	earlier := &v486.StartGame{
		EntityUniqueID:                 latest.EntityUniqueID,
		EntityRuntimeID:                latest.EntityRuntimeID,
		PlayerGameMode:                 latest.PlayerGameMode,
		PlayerPosition:                 latest.PlayerPosition,
		Pitch:                          latest.Pitch,
		Yaw:                            latest.Yaw,
		WorldSeed:                      int32(latest.WorldSeed),
		SpawnBiomeType:                 latest.SpawnBiomeType,
		UserDefinedBiomeName:           latest.UserDefinedBiomeName,
		Dimension:                      latest.Dimension,
		Generator:                      latest.Generator,
		WorldGameMode:                  latest.WorldGameMode,
		Difficulty:                     latest.Difficulty,
		WorldSpawn:                     latest.WorldSpawn,
		AchievementsDisabled:           latest.AchievementsDisabled,
		DayCycleLockTime:               latest.DayCycleLockTime,
		EducationEditionOffer:          latest.EducationEditionOffer,
		EducationFeaturesEnabled:       latest.EducationFeaturesEnabled,
		EducationProductID:             latest.EducationProductID,
		RainLevel:                      latest.RainLevel,
		LightningLevel:                 latest.LightningLevel,
		ConfirmedPlatformLockedContent: latest.ConfirmedPlatformLockedContent,
		MultiPlayerGame:                latest.MultiPlayerGame,
		MultiPlayerCorrelationID:       latest.MultiPlayerCorrelationID,
		LANBroadcastEnabled:            latest.LANBroadcastEnabled,
		XBLBroadcastMode:               latest.XBLBroadcastMode,
		CommandsEnabled:                latest.CommandsEnabled,
		TexturePackRequired:            latest.TexturePackRequired,
		GameRules:                      latest.GameRules,
		Experiments:                    latest.Experiments,
		ExperimentsPreviouslyToggled:   latest.ExperimentsPreviouslyToggled,
		BonusChestEnabled:              latest.BonusChestEnabled,
		StartWithMapEnabled:            latest.StartWithMapEnabled,
		PlayerPermissions:              latest.PlayerPermissions,
		ServerChunkTickRadius:          latest.ServerChunkTickRadius,
		HasLockedBehaviourPack:         latest.HasLockedBehaviourPack,
		HasLockedTexturePack:           latest.HasLockedTexturePack,
		FromLockedWorldTemplate:        latest.FromLockedWorldTemplate,
		MSAGamerTagsOnly:               latest.MSAGamerTagsOnly,
		FromWorldTemplate:              latest.FromWorldTemplate,
		WorldTemplateSettingsLocked:    latest.WorldTemplateSettingsLocked,
		OnlySpawnV1Villagers:           latest.OnlySpawnV1Villagers,
		BaseGameVersion:                latest.BaseGameVersion,
		LimitedWorldWidth:              latest.LimitedWorldWidth,
		LimitedWorldDepth:              latest.LimitedWorldDepth,
		NewNether:                      latest.NewNether,
		EducationSharedResourceURI:     latest.EducationSharedResourceURI,
		ForceExperimentalGameplay:      latest.ForceExperimentalGameplay,
		LevelID:                        latest.LevelID,
		WorldName:                      latest.WorldName,
		TemplateContentIdentity:        latest.TemplateContentIdentity,
		Trial:                          latest.Trial,
		PlayerMovementSettings:         latest.PlayerMovementSettings,
		Time:                           latest.Time,
		EnchantmentSeed:                latest.EnchantmentSeed,
		Blocks:                         latest.Blocks,
		ServerAuthoritativeInventory:   latest.ServerAuthoritativeInventory,
		GameVersion:                    latest.GameVersion,
		ServerBlockStateChecksum:       latest.ServerBlockStateChecksum,
	}
	for _, i := range latest.Items {
		if oldRuntimeID, ok := legacymappings.ItemNameToRuntimeID(i.Name); ok {
			earlier.Items = append(earlier.Items, protocol.ItemEntry{
				Name:           i.Name,
				RuntimeID:      int16(oldRuntimeID),
				ComponentBased: i.ComponentBased,
			})
		}
	}
	return earlier
}