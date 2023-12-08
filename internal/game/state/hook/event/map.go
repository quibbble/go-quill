package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

var EventMap map[string]func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error

func init() {
	EventMap = map[string]func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error{
		DamageUnitEvent:         DamageUnitAffect,
		DamageUnitsEvent:        DamageUnitsAffect,
		KillUnitEvent:           KillUnitAffect,
		RescindUnitEvent:        RescindUnitAffect,
		HealUnitEvent:           HealUnitAffect,
		HealUnitsEvent:          HealUnitsAffect,
		MoveUnitEvent:           MoveUnitAffect,
		AttackUnitEvent:         AttackUnitAffect,
		PlaceUnitEvent:          PlaceUnitAffect,
		SwapUnitsEvent:          SwapUnitsAffect,
		SummonUnitEvent:         SummonUnitAffect,
		ModifyUnitEvent:         ModifyUnitAffect,
		AddItemToUnitEvent:      AddItemToUnitAffect,
		RemoveItemFromUnitEvent: RemoveItemFromUnitAffect,

		RefreshMovementEvent: RefreshMovementAffect,
		CooldownEvent:        CooldownAffect,

		PlayCardEvent:       PlayCardAffect,
		DiscardCardEvent:    DiscardCardAffect,
		DrawCardEvent:       DrawCardAffect,
		BurnCardEvent:       BurnCardAffect,
		SackCardEvent:       SackCardAffect,
		AddTraitToCard:      AddTraitToCardAffect,
		RemoveTraitFromCard: RemoveTraitFromCardAffect,

		RecycleDeckEvent: RecycleDeckAffect,

		DrainManaEvent:     DrainManaAffect,
		GainManaEvent:      GainManaAffect,
		DrainBaseManaEvent: DrainBaseManaAffect,
		GainBaseManaEvent:  GainBaseManaAffect,

		EndTurnEvent: EndTurnAffect,
	}
}
