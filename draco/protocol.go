package draco

import (
	"bytes"
	"fmt"

	"github.com/cqdetdev/draco/draco/legacy"
	"github.com/cqdetdev/draco/draco/state"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// protcool representst he
type Protocol struct {
	minecraft.Protocol
}

func (Protocol) ID() int32 {
	return 486
}

func (Protocol) Ver() string {
	return "1.18.10"
}

func (p Protocol) Packets() packet.Pool { return packet.NewPool() }
func (p Protocol) ConvertToLatest(pk packet.Packet) packet.Packet {
	switch pk.ID() {
	// debug only
	case packet.IDResourcePackClientResponse:
		r := pk.(*packet.ResourcePackClientResponse)
		fmt.Printf("Response %d\n", r.Response)
		fmt.Printf("Packs: %v\n", r.PacksToDownload)
		return r
	default:
		return pk
	}
}
func (p Protocol) ConvertFromLatest(pk packet.Packet) packet.Packet {
	switch pk.ID() {
	case packet.IDCreativeContent:
		c := pk.(*packet.CreativeContent)
		// iterate over items
		for i := 0; i < len(c.Items); i++ {
		}
	case packet.IDPacketViolationWarning:
		v := pk.(*packet.PacketViolationWarning)
		fmt.Printf("Violation %d\n", v.PacketID)
		return v
	case packet.IDAddPlayer:
		old := pk.(*packet.AddPlayer)
		new := &legacy.AddPlayer{}
		new.UUID = old.UUID
		new.Username = old.Username
		new.EntityUniqueID = old.EntityUniqueID
		new.EntityRuntimeID = old.EntityRuntimeID
		new.PlatformChatID = old.PlatformChatID
		new.Position = old.Position
		new.Velocity = old.Velocity
		new.Pitch = old.Pitch
		new.Yaw = old.Yaw
		new.HeadYaw = old.HeadYaw
		old.HeldItem.Stack.NetworkID = 0
		new.HeldItem = old.HeldItem
		old.HeldItem.StackNetworkID = 0
		if old.HeldItem.Stack.HasNetworkID {
			new.HeldItem.StackNetworkID = old.HeldItem.StackNetworkID
		}
		new.EntityMetadata = old.EntityMetadata
		new.Flags = old.Flags
		new.CommandPermissionLevel = old.CommandPermissionLevel
		new.ActionPermissions = old.ActionPermissions
		new.PermissionLevel = old.PermissionLevel
		new.CustomStoredPermissions = old.CustomStoredPermissions
		new.PlayerUniqueID = old.PlayerUniqueID
		new.EntityLinks = old.EntityLinks
		new.DeviceID = old.DeviceID
		new.BuildPlatform = old.BuildPlatform
		return new
	case packet.IDStartGame:
		old := pk.(*packet.StartGame)
		new := &legacy.StartGame{}
		new.EntityUniqueID = old.EntityUniqueID
		new.EntityRuntimeID = old.EntityRuntimeID
		new.PlayerGameMode = old.PlayerGameMode
		new.PlayerPosition = old.PlayerPosition
		new.Pitch = old.Pitch
		new.Yaw = old.Yaw
		new.WorldSeed = int32(old.WorldSeed)
		new.SpawnBiomeType = old.SpawnBiomeType
		new.UserDefinedBiomeName = old.UserDefinedBiomeName
		new.Dimension = old.Dimension
		new.Generator = old.Generator
		new.WorldGameMode = old.WorldGameMode
		new.Difficulty = old.Difficulty
		new.WorldSpawn = old.WorldSpawn
		new.AchievementsDisabled = old.AchievementsDisabled
		new.DayCycleLockTime = old.DayCycleLockTime
		new.EducationEditionOffer = old.EducationEditionOffer
		new.EducationFeaturesEnabled = old.EducationFeaturesEnabled
		new.EducationProductID = old.EducationProductID
		new.RainLevel = old.RainLevel
		new.LightningLevel = old.LightningLevel
		new.ConfirmedPlatformLockedContent = old.ConfirmedPlatformLockedContent
		new.MultiPlayerGame = old.MultiPlayerGame
		new.MultiPlayerCorrelationID = old.MultiPlayerCorrelationID
		new.LANBroadcastEnabled = old.LANBroadcastEnabled
		new.XBLBroadcastMode = old.XBLBroadcastMode
		new.CommandsEnabled = old.CommandsEnabled
		new.TexturePackRequired = old.TexturePackRequired
		new.GameRules = old.GameRules
		new.Experiments = old.Experiments
		new.ExperimentsPreviouslyToggled = old.ExperimentsPreviouslyToggled
		new.BonusChestEnabled = old.BonusChestEnabled
		new.StartWithMapEnabled = old.StartWithMapEnabled
		new.PlayerPermissions = old.PlayerPermissions
		new.ServerChunkTickRadius = old.ServerChunkTickRadius
		new.HasLockedBehaviourPack = old.HasLockedBehaviourPack
		new.HasLockedTexturePack = old.HasLockedTexturePack
		new.FromLockedWorldTemplate = old.FromLockedWorldTemplate
		new.MSAGamerTagsOnly = old.MSAGamerTagsOnly
		new.FromWorldTemplate = old.FromWorldTemplate
		new.WorldTemplateSettingsLocked = old.WorldTemplateSettingsLocked
		new.OnlySpawnV1Villagers = old.OnlySpawnV1Villagers
		new.BaseGameVersion = old.BaseGameVersion
		new.LimitedWorldWidth = old.LimitedWorldWidth
		new.LimitedWorldDepth = old.LimitedWorldDepth
		new.NewNether = old.NewNether
		new.EducationSharedResourceURI = old.EducationSharedResourceURI
		new.ForceExperimentalGameplay = old.ForceExperimentalGameplay
		new.LevelID = old.LevelID
		new.WorldName = old.WorldName
		new.TemplateContentIdentity = old.TemplateContentIdentity
		new.Trial = old.Trial
		new.PlayerMovementSettings = old.PlayerMovementSettings
		new.Time = old.Time
		new.EnchantmentSeed = old.EnchantmentSeed
		new.Blocks = old.Blocks
		for _, i := range old.Items {
			it, ok := world.ItemByRuntimeID(int32(i.RuntimeID), 0)
			name, _ := it.EncodeItem()
			if !ok {
				panic("could not convert it idk (to new)")
			}

			oldRuntimeID, meta, ok := state.ItemRuntimeID(it)
			if !ok {
				fmt.Printf("%s %v\n", name, meta)
				panic("could not convert it idk (to old)")
			}
			fmt.Printf("(%s) New: %d | Old: %d\n", name, i.RuntimeID, oldRuntimeID)
			new.Items = append(new.Items, protocol.ItemEntry{
				Name:           name,
				RuntimeID:      int16(oldRuntimeID),
				ComponentBased: i.ComponentBased,
			})
		}
		new.MultiPlayerCorrelationID = old.MultiPlayerCorrelationID
		new.ServerAuthoritativeInventory = old.ServerAuthoritativeInventory
		new.GameVersion = old.GameVersion
		new.ServerBlockStateChecksum = old.ServerBlockStateChecksum
		return new
	case packet.IDLevelChunk:
		lc := pk.(*packet.LevelChunk)
		air := world.BlockRuntimeID(block.Air{})
		c, err := chunk.NetworkDecode(air, lc.RawPayload, int(lc.SubChunkCount), world.Overworld.Range())
		if err != nil {
			panic(err)
		}
		for _, s := range c.Sub() {
			for _, l := range s.Layers() {
				l.Palette().Replace(func(newRuntimeID uint32) uint32 {
					name, props, ok := chunk.RuntimeIDToState(newRuntimeID)
					if !ok {
						panic("could not convert it idk (to new)")
					}

					oldRuntimeID, ok := state.StateToRuntimeID(name, props)
					if !ok {
						fmt.Printf("%s %v\n", name, props)
						panic("could not convert it idk (to old)")
					}
					// fmt.Printf("(%s) New: %d | Old: %d\n", name, newRuntimeID, oldRuntimeID)
					return oldRuntimeID
				})
			}
		}

		buf, data := bytes.NewBuffer(nil), chunk.Encode(c, chunk.NetworkEncoding)
		for i := range data.SubChunks {
			_, _ = buf.Write(data.SubChunks[i])
		}
		_, _ = buf.Write(data.Biomes)
		buf.WriteByte(0)
		lc.RawPayload = buf.Bytes()
		return lc
	case packet.IDAddVolumeEntity:
		old := pk.(*packet.AddVolumeEntity)
		new := &legacy.AddVolumeEntity{}
		new.EntityRuntimeID = old.EntityRuntimeID
		new.EntityMetadata = old.EntityMetadata
		new.EncodingIdentifier = old.EncodingIdentifier
		new.InstanceIdentifier = old.InstanceIdentifier
		new.EngineVersion = old.EngineVersion
		return new
	case packet.IDPlayStatus:
		fmt.Println("PlayStatus")
		return pk
	}

	return pk
}
