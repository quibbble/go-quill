package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	dg "github.com/quibbble/go-quill/internal/state/damage"
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
	Choose     Choose
}

func DamageUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DamageUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), a.Choose.Type, a.Choose.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
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
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
	damage, err := dg.Damage(unit, a.Amount, a.DamageType)
	if err != nil {
		return errors.Wrap(err)
	}
	unit.Health -= damage

	if unit.Health <= 0 {
		event := &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  KillUnitEvent,
			args: &KillUnitArgs{
				Choose: Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: unit.UUID,
					},
				},
			},
			affect: KillUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	} else {
		// enrage trait check
		for _, trait := range unit.GetTraits(tr.EnrageTrait) {
			args := trait.GetArgs().(*tr.EnrageArgs)
			event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
			if err != nil {
				return errors.Wrap(err)
			}
			if err := engine.Do(event, state); err != nil {
				return errors.Wrap(err)
			}
		}
	}

	return nil
}
