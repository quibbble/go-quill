package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	dg "github.com/quibbble/go-quill/internal/state/damage"
	"github.com/quibbble/go-quill/internal/state/hook/choose"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DamageUnitEvent = "DamageUnit"
)

type DamageUnitArgs struct {
	DamageType string
	Amount     int
	ch.Choose
}

func DamageUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(DamageUnitArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	choices, err := a.Choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	if len(choices) != 1 {
		return errors.ErrInvalidSliceLength
	}
	if choices[0].Type() != st.UnitUUID {
		return st.ErrInvalidUUIDType(choices[0], st.UnitUUID)
	}
	x, y, err := state.Board.GetUnitXY(choices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit
	damage, err := dg.Damage(unit, a.Amount, a.DamageType)
	if err != nil {
		return errors.Wrap(err)
	}
	unit.Health -= damage
	if unit.Health <= 0 {
		event := &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  KillUnitEvent,
			args: &KillUnitArgs{
				&choose.UUIDChoice{
					UUID: unit.UUID,
				},
			},
			affect: KillUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
