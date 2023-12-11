package event

import (
	"context"

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

func DamageUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DamageUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
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
		event, err := NewEvent(state.Gen.New(en.EventUUID), KillUnitEvent, KillUnitArgs{
			ChooseUnit: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: unit.UUID,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	} else {
		// enrage trait check
		for _, trait := range unit.GetTraits(tr.EnrageTrait) {
			args := trait.GetArgs().(*tr.EnrageArgs)
			for _, h := range args.Hooks {
				hook, err := state.NewHook(state.Gen, unit.GetUUID(), h)
				if err != nil {
					return errors.Wrap(err)
				}
				engine.Register(hook)
			}
			for _, e := range args.Events {
				event, err := NewEvent(state.Gen.New(en.EventUUID), e.Type, e.Args)
				if err != nil {
					return errors.Wrap(err)
				}
				ctx := context.WithValue(context.WithValue(context.Background(), en.TraitCardCtx, unit.GetUUID()), en.TraitEventCtx, e)
				if err := engine.Do(ctx, event, state); err != nil {
					return errors.Wrap(err)
				}
			}
		}
	}

	return nil
}
