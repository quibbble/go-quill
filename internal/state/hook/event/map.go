package event

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var (
	EventMap = map[string]func(engine *engine.Engine, state *state.State, args interface{}, targets ...uuid.UUID) error{
		DamageUnitEvent:  DamageUnitAffect,
		DamageUnitsEvent: DamageUnitsAffect,
		KillUnitEvent:    KillUnitAffect,
		HealUnitEvent:    HealUnitAffect,
		MoveUnitEvent:    MoveUnitAffect,
		AttackUnitEvent:  AttackUnitAffect,
		PlaceUnitEvent:   PlaceUnitAffect,
		SummonUnitEvent:  SummonUnitAffect,
	}
)
