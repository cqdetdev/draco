package draco

import (
	"bytes"
	"fmt"

	"github.com/cqdetdev/draco/draco/chunk"
	"github.com/cqdetdev/draco/draco/latestmappings"
	"github.com/cqdetdev/draco/draco/legacymappings"
	"github.com/cqdetdev/draco/draco/legacypackets"
	"github.com/cqdetdev/draco/draco/legacystructures"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
)

// currentFormID is the ID of the current form that the player has open
var currentFormID uint32

var sentModalForm bool

var changedCommandArgs = []uint32{
	7, 9, 16, 38, 46, 47, 50, 52, 56, 69,
}

// Protocol is the protocol used to support the Minecraft 1.18.10 protocol (486).
type Protocol struct {
	minecraft.Protocol
}

// ID ...
func (Protocol) ID() int32 {
	return 486
}

// Ver ...
func (Protocol) Ver() string {
	return "1.18.12"
}

// Packets ...
func (p Protocol) Packets() packet.Pool {
	pool := packet.NewPool()
	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &legacypackets.PlayerAuthInput{} }
	pool[packet.IDModalFormResponse] = func() packet.Packet { return &legacypackets.ModalFormResponse{} }
	return pool
}

var (
	// air is the runtime ID of an air block.
	air, _ = latestmappings.StateToRuntimeID("minecraft:air", nil)
	// infoBlock is the runtime ID of an info block.
	infoBlock, _ = latestmappings.StateToRuntimeID("minecraft:info_update", nil)
	// infoItem is the runtime ID of an info item.
	infoItem, _ = latestmappings.ItemNameToRuntimeID("minecraft:info_update")
	// worldRange is hardcoded to the overworld world range.
	// TODO: Dimensions support.
	worldRange = cube.Range{-64, 319}
)

// ConvertToLatest ...
func (p Protocol) ConvertToLatest(pk packet.Packet, c *minecraft.Conn) []packet.Packet {
	switch earlier := pk.(type) {
	case *packet.PacketViolationWarning:
		logrus.Infof("Violation %X (%d): %v\n", earlier.PacketID, earlier.Severity, earlier.ViolationContext)
	case *packet.MobEquipment:
		earlier.NewItem.Stack = upgradeItemStack(earlier.NewItem.Stack)
		return []packet.Packet{earlier}
	case *legacypackets.PlayerAuthInput:
		latest := &packet.PlayerAuthInput{}
		latest.Pitch = earlier.Pitch
		latest.Yaw = earlier.Yaw
		latest.Position = earlier.Position
		latest.HeadYaw = earlier.HeadYaw
		latest.MoveVector = earlier.MoveVector
		latest.InputData = earlier.InputData
		latest.InputMode = earlier.InputMode
		latest.PlayMode = earlier.PlayMode
		latest.InteractionModel = packet.InteractionModelClassic // ?
		latest.GazeDirection = earlier.GazeDirection
		latest.Tick = earlier.Tick
		latest.Delta = earlier.Delta
		latest.ItemInteractionData = earlier.ItemInteractionData
		latest.ItemStackRequest = earlier.ItemStackRequest
		latest.BlockActions = earlier.BlockActions
		return []packet.Packet{latest}
	case *legacypackets.ModalFormResponse:
		if sentModalForm {
			sentModalForm = false
			break
		}
		latest := &packet.ModalFormResponse{}
		latest.FormID = currentFormID
		latest.ResponseData = protocol.Option(earlier.ResponseData)
		latest.CancelReason = protocol.Option[uint8](packet.ModalFormCancelReasonUserClosed)
		if !sentModalForm {
			sentModalForm = true
		}
		return []packet.Packet{latest}
	case *packet.InventoryTransaction:
		actions := make([]protocol.InventoryAction, 0, len(earlier.Actions))
		for _, action := range earlier.Actions {
			action.OldItem.Stack = upgradeItemStack(action.OldItem.Stack)
			action.NewItem.Stack = upgradeItemStack(action.NewItem.Stack)
			actions = append(actions, action)
		}
		earlier.Actions = actions
		switch data := earlier.TransactionData.(type) {
		case *protocol.UseItemTransactionData:
			data.HeldItem.Stack = upgradeItemStack(data.HeldItem.Stack)
			data.BlockRuntimeID = upgradeBlockRuntimeID(data.BlockRuntimeID)
		case *protocol.UseItemOnEntityTransactionData:
			data.HeldItem.Stack = upgradeItemStack(data.HeldItem.Stack)
		}
	default:
		return []packet.Packet{pk}
	}

	return nil
}

// ConvertFromLatest ...
func (p Protocol) ConvertFromLatest(pk packet.Packet, _ *minecraft.Conn) []packet.Packet {
	switch latest := pk.(type) {
	case *packet.PacketViolationWarning:
		logrus.Infof("Violation %X (%d): (Context: %s)\n", latest.PacketID, latest.Severity, latest.ViolationContext)
	case *packet.UpdateBlock:
		latest.NewBlockRuntimeID = downgradeBlockRuntimeID(latest.NewBlockRuntimeID)
	case *packet.SetActorData:
		downgradeEntityMetadata(latest.EntityMetadata)
	case *packet.ModalFormRequest:
		currentFormID = latest.FormID
	case *packet.AvailableCommands:
		// cmds := latest.Commands
		// TODO: Update command arg values...
		return []packet.Packet{latest}
	case *packet.AddActor:
		earlier := &legacypackets.AddActor{}
		earlier.EntityUniqueID = latest.EntityUniqueID
		earlier.EntityRuntimeID = latest.EntityRuntimeID
		earlier.EntityType = latest.EntityType
		earlier.Position = latest.Position
		earlier.Velocity = latest.Velocity
		earlier.Pitch = latest.Pitch
		earlier.Yaw = latest.Yaw
		earlier.HeadYaw = latest.HeadYaw
		var earlierAttributes []protocol.Attribute
		for _, attr := range latest.Attributes {
			earlierAttributes = append(earlierAttributes, protocol.Attribute{
				AttributeValue: attr,
			})
		}
		earlier.Attributes = earlierAttributes
		downgradeEntityMetadata(latest.EntityMetadata)
		earlier.EntityMetadata = latest.EntityMetadata
		return []packet.Packet{earlier}
	case *packet.CraftingData:
		recipes := make([]protocol.Recipe, 0, len(latest.Recipes))
		for _, r := range latest.Recipes {
			switch r := r.(type) {
			case *protocol.ShapedRecipe:
				r.Input, r.Output = downgradeRecipe(r.Input, r.Output)
			case *protocol.ShapelessRecipe:
				r.Input, r.Output = downgradeRecipe(r.Input, r.Output)
			}
			recipes = append(recipes, r)
		}
		latest.Recipes = recipes
	case *packet.CreativeContent:
		items := make([]protocol.CreativeItem, 0, len(latest.Items))
		for _, it := range latest.Items {
			it.Item = downgradeItemStack(it.Item)
			items = append(items, it)
		}
		latest.Items = items
	case *packet.InventoryContent:
		items := make([]protocol.ItemInstance, 0, len(latest.Content))
		for _, it := range latest.Content {
			it.Stack = downgradeItemStack(it.Stack)
			items = append(items, it)
		}
		latest.Content = items
	case *packet.InventorySlot:
		latest.NewItem.Stack = downgradeItemStack(latest.NewItem.Stack)
	case *packet.UpdateAttributes:
		earlier := &legacypackets.UpdateAttributes{}
		earlier.EntityRuntimeID = latest.EntityRuntimeID
		earlier.Attributes = []legacystructures.Attribute{}
		for _, attr := range latest.Attributes {
			earlier.Attributes = append(earlier.Attributes, legacystructures.Attribute{
				Name:    attr.Name,
				Value:   attr.Value,
				Max:     attr.Max,
				Min:     attr.Min,
				Default: attr.Default,
			})
		}
		return []packet.Packet{earlier}
	case *packet.NetworkChunkPublisherUpdate:
		earlier := &legacypackets.NetworkChunkPublisherUpdate{}
		earlier.Position = latest.Position
		earlier.Radius = latest.Radius
		return []packet.Packet{earlier}
	case *packet.AddPlayer:
		earlier := &legacypackets.AddPlayer{
			UUID:            latest.UUID,
			Username:        latest.Username,
			EntityUniqueID:  latest.EntityUniqueID,
			EntityRuntimeID: latest.EntityRuntimeID,
			PlatformChatID:  latest.PlatformChatID,
			Position:        latest.Position,
			Velocity:        latest.Velocity,
			Pitch:           latest.Pitch,
			Yaw:             latest.Yaw,
			HeadYaw:         latest.HeadYaw,
			HeldItem:        latest.HeldItem,
			EntityMetadata:  latest.EntityMetadata,
			EntityLinks:     latest.EntityLinks,
			DeviceID:        latest.DeviceID,
			BuildPlatform:   latest.BuildPlatform,
		}
		earlier.HeldItem.Stack = downgradeItemStack(latest.HeldItem.Stack)
		return []packet.Packet{earlier}
	case *packet.StartGame:
		feg, _ := latest.ForceExperimentalGameplay.Value()
		earlier := &legacypackets.StartGame{
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
			PlayerPermissions:              int32(latest.PlayerPermissions),
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
			ForceExperimentalGameplay:      feg,
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
		return []packet.Packet{earlier}
	case *packet.LevelChunk:
		if latest.SubChunkRequestMode == protocol.SubChunkRequestModeLegacy {
			readBuf := bytes.NewBuffer(latest.RawPayload)
			c, err := chunk.NetworkDecode(air, readBuf, int(latest.SubChunkCount), worldRange)
			if err != nil {
				panic(err)
			}
			for _, s := range c.Sub() {
				downgradeSubChunk(s)
			}

			writeBuf, data := bytes.NewBuffer(nil), chunk.Encode(c, chunk.NetworkEncoding)
			for i := range data.SubChunks {
				_, _ = writeBuf.Write(data.SubChunks[i])
			}
			_, _ = writeBuf.Write(data.Biomes)

			latest.RawPayload = append(writeBuf.Bytes(), readBuf.Bytes()...)
		}
	case *packet.SubChunk:
		entries := make([]protocol.SubChunkEntry, 0, len(latest.SubChunkEntries))
		for _, e := range latest.SubChunkEntries {
			if e.Result == protocol.SubChunkResultSuccess {
				var ind uint8
				buf := bytes.NewBuffer(e.RawPayload)
				s, err := chunk.DecodeSubChunk(air, worldRange, buf, &ind, chunk.NetworkEncoding)
				if err != nil {
					panic(err)
				}
				downgradeSubChunk(s)
				serialisedSubChunk := chunk.EncodeSubChunk(s, chunk.NetworkEncoding, worldRange, int(ind))
				e.RawPayload = append(serialisedSubChunk, buf.Bytes()...)
			}
			entries = append(entries, e)
		}
		latest.SubChunkEntries = entries
	case *packet.AddVolumeEntity:
		return []packet.Packet{&legacypackets.AddVolumeEntity{
			EntityRuntimeID:    latest.EntityRuntimeID,
			EntityMetadata:     latest.EntityMetadata,
			EncodingIdentifier: latest.EncodingIdentifier,
			InstanceIdentifier: latest.InstanceIdentifier,
			EngineVersion:      latest.EngineVersion,
		}}
	case *packet.RemoveVolumeEntity:
		return []packet.Packet{&legacypackets.RemoveVolumeEntity{EntityRuntimeID: latest.EntityRuntimeID}}
	case *packet.SpawnParticleEffect:
		return []packet.Packet{&legacypackets.SpawnParticleEffect{
			Dimension:      latest.Dimension,
			EntityUniqueID: latest.EntityUniqueID,
			Position:       latest.Position,
			ParticleName:   latest.ParticleName,
		}}
	default:
		return []packet.Packet{pk}
	}
	return nil
}

// dataKeyVariant is used for falling blocks and fake texts. This is necessary for falling block runtime ID translation.
const dataKeyVariant = 2

// downgradeSubChunk translates a 1.19.10 sub-chunk to a 1.18.12 one, updating all palette entries with the appropriate
// runtime IDs.
func downgradeSubChunk(s *chunk.SubChunk) {
	for _, l := range s.Layers() {
		l.Palette().Replace(downgradeBlockRuntimeID)
	}
}

// downgradeBlockRuntimeID translates a 1.19.10 runtime ID to a 1.18.12 one.
func downgradeBlockRuntimeID(latestRID uint32) uint32 {
	name, properties, found := latestmappings.RuntimeIDToState(latestRID)
	if !found {
		panic(fmt.Errorf("downgrade block runtime id: could not find name for runtime id: %v", latestRID))
	}
	earlierRuntimeID, found := legacymappings.StateToRuntimeID(name, properties)
	if !found {
		// logrus.Errorf("downgrade block runtime id: could not find runtime id for name: %v", name)
		return infoBlock
	}
	return earlierRuntimeID
}

// upgradeBlockRuntimeID translates a 1.18.12 block runtime ID to a 1.19.10 one.
func upgradeBlockRuntimeID(id uint32) uint32 {
	name, properties, found := legacymappings.RuntimeIDToState(id)
	if !found {
		panic(fmt.Errorf("upgrade block runtime id: could not find name for runtime id: %v", id))
	}
	latestRuntimeID, found := latestmappings.StateToRuntimeID(name, properties)
	if !found {
		panic(fmt.Errorf("upgrade block runtime id: could not find runtime id for name: %v", name))
	}
	return latestRuntimeID
}

// downgradeEntityMetadata translates a 1.19.10 entity metadata to a 1.18.12 one.
func downgradeEntityMetadata(metadata map[uint32]any) {
	if latestRID, ok := metadata[dataKeyVariant]; ok {
		metadata[dataKeyVariant] = int32(downgradeBlockRuntimeID(uint32(latestRID.(int32))))
	}
}

// downgradeRecipe downgrades a 1.19.10 recipe to a 1.18.12 one.
func downgradeRecipe(latestInput []protocol.RecipeIngredientItem, latestOutput []protocol.ItemStack) ([]protocol.RecipeIngredientItem, []protocol.ItemStack) {
	input := make([]protocol.RecipeIngredientItem, 0, len(latestInput))
	output := make([]protocol.ItemStack, 0, len(latestOutput))
	for _, i := range latestInput {
		if i.Count > 0 {
			i.NetworkID = downgradeItemRuntimeID(i.NetworkID)
		}
		input = append(input, i)
	}
	for _, o := range latestOutput {
		output = append(output, downgradeItemStack(o))
	}
	return input, output
}

// downgradeItemStack translates a 1.19.10 item stack to a 1.18.12 one, updating all palette entries with the appropriate
// runtime IDs.
func downgradeItemStack(st protocol.ItemStack) protocol.ItemStack {
	if st.BlockRuntimeID > 0 {
		st.BlockRuntimeID = int32(downgradeBlockRuntimeID(uint32(st.BlockRuntimeID)))
	}
	if st.HasNetworkID {
		st.NetworkID = downgradeItemRuntimeID(st.NetworkID)
	}
	return st
}

// upgradeItemStack translates a 1.18.12 item stack to a 1.19.10 one, updating all palette entries with the appropriate
// runtime IDs.
func upgradeItemStack(st protocol.ItemStack) protocol.ItemStack {
	if st.BlockRuntimeID > 0 {
		st.BlockRuntimeID = int32(upgradeBlockRuntimeID(uint32(st.BlockRuntimeID)))
	}
	if st.HasNetworkID {
		st.NetworkID = upgradeItemRuntimeID(st.NetworkID)
	}
	return st
}

// upgradeItemRuntimeID translates a 1.18.12 item runtime ID to a 1.19.10 one.
func upgradeItemRuntimeID(latestRID int32) int32 {
	name, found := legacymappings.ItemRuntimeIDToName(latestRID)
	if !found {
		panic(fmt.Errorf("upgrade item runtime id: could not find name for runtime id: %v", latestRID))
	}
	earlierRuntimeID, found := latestmappings.ItemNameToRuntimeID(name)
	if !found {
		panic(fmt.Errorf("upgrade item runtime id: could not find runtime id for name: %v", name))
	}
	return earlierRuntimeID
}

// downgradeItemRuntimeID translates a 1.19.10 item runtime ID to a 1.18.12 one.
func downgradeItemRuntimeID(latestRID int32) int32 {
	name, found := latestmappings.ItemRuntimeIDToName(latestRID)
	if !found {
		panic(fmt.Errorf("downgrade item runtime id: could not find name for runtime id: %v", latestRID))
	}
	earlierRuntimeID, found := legacymappings.ItemNameToRuntimeID(name)
	if !found {
		// logrus.Errorf("downgrade item runtime id: could not find runtime id for name: %v", name)
		return infoItem
	}
	return earlierRuntimeID
}
