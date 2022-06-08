package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var Translator translator.Translator = translator.Translator{
	Inbound: map[uint32]translator.TranslationHandler{
		packet.IDUpdateBlock: UpdateBlockTranslator{},
		packet.IDSetActorData: SetActorDataTranslator{},
		packet.IDAddActor: AddActorTranslator{},
		packet.IDCraftingData: CraftingDataTranslator{},
		packet.IDCreativeContent: CreativeContentTranslator{},
		packet.IDInventoryContent: InventoryContentTranslator{},
		packet.IDInventorySlot: InventorySlotTranslator{},
		packet.IDAddPlayer: AddPlayerTranslator{},
		packet.IDStartGame: StartGameTranslator{},
		packet.IDLevelChunk: LevelChunkTranslator{},
		packet.IDSubChunk: SubChunkTranslator{},
		packet.IDAddVolumeEntity: AddVolumeEntityTranslator{},
		packet.IDRemoveVolumeEntity: RemoveVolumeEntityTranslator{},
		packet.IDSpawnParticleEffect: SpawnParticleEffectTranslator{},
		packet.IDSetTitle: SetTitleTranslator{},
	},
	Outbound: map[uint32]translator.TranslationHandler{
		packet.IDMobEquipment: MobEquipmentTranslator{},
		packet.IDPlayerAuthInput: PlayerAuthInputTranslator{},
		packet.IDInventoryTransaction: InventoryTransactionTranslator{},
		packet.IDResourcePacksInfo: ResourcePacksInfoTranslator{},
	},
	Protocol: 486,
}