package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/internal/state/hook/choose"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DamageUnitsEvent = "damage_units"
)

type DamageUnitsArgs struct {
	DamageType string
	Amount     int
	ch.Choose
}

func DamageUnitsAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(DamageUnitArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	choices, err := a.Choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		if choice.Type() != st.UnitUUID {
			return st.ErrInvalidUUIDType(choice, st.UnitUUID)
		}
		// it's possible the unit was killed by another affect between choosing and applying
		// if this happens then just continue
		_, _, err := state.Board.GetUnitXY(choice)
		if err != nil {
			continue
		}
		event := &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DamageUnitEvent,
			args: &DamageUnitArgs{
				DamageType: a.DamageType,
				Amount:     a.Amount,
				Choose: &choose.UUIDChoice{
					UUID: choice,
				},
			},
			affect: DamageUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
