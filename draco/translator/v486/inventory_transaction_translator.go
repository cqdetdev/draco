package v486

import (
	"github.com/cqdetdev/draco/draco/translator"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type InventoryTransactionTranslator struct{}

func (InventoryTransactionTranslator) Translate(pk packet.Packet) packet.Packet {
	latest := pk.(*packet.InventoryTransaction)
	actions := make([]protocol.InventoryAction, 0, len(latest.Actions))
	for _, action := range latest.Actions {
		action.OldItem.Stack = translator.UpgradeItemStack(action.OldItem.Stack)
		action.NewItem.Stack = translator.UpgradeItemStack(action.NewItem.Stack)
		actions = append(actions, action)
	}
	latest.Actions = actions
	switch data := latest.TransactionData.(type) {
	case *protocol.UseItemTransactionData:
		data.HeldItem.Stack = translator.UpgradeItemStack(data.HeldItem.Stack)
		data.BlockRuntimeID = translator.UpgradeBlockRuntimeID(data.BlockRuntimeID)
	case *protocol.UseItemOnEntityTransactionData:
		data.HeldItem.Stack = translator.UpgradeItemStack(data.HeldItem.Stack)
	}
	return latest
}