package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	dg "github.com/quibbble/go-quill/internal/game/state/damage"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	DamageUnitEvent = "DamageUnit"
)

type DamageUnitArgs struct {
	DamageType string
	Amount     int
	ChooseUnit parse.Choose
}

func DamageUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a DamageUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	unitChoice, err := GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
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
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: unit.UUID,
					},
				},
			},
			affect: KillUnitAffect,
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	} else {
		// enrage trait check
		for _, trait := range unit.GetTraits(tr.EnrageTrait) {
			args := trait.GetArgs().(tr.EnrageArgs)
			event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
			if err != nil {
				return errors.Wrap(err)
			}
			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}
		}
	}

	return nil
}
