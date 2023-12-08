package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	KillUnitEvent = "KillUnit"
)

type KillUnitArgs struct {
	ChooseUnit parse.Choose
}

func KillUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a KillUnitArgs
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
	state.Board.XYs[x][y].Unit = nil

	// death cry trait check
	deathCrys := unit.GetTraits(tr.DeathCryTrait)
	for _, deathCry := range deathCrys {
		args := deathCry.GetArgs().(tr.DeathCryArgs)
		event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	if unit.GetID() != "U0001" {
		// reset and move items and unit to discard
		for _, item := range unit.Items {
			if item.Player == unit.Player {
				item.Reset(state.BuildCard)
				state.Discard[unit.Player].Add(item)
			}
		}
		unit.Reset(state.BuildCard)
		state.Discard[unit.Player].Add(unit)
	} else {
		// check if the game is over
		choose1, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.UnitsChoice, ch.UnitsArgs{
			Types: []string{cd.BaseUnit},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		choose2, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.OwnedUnitsChoice, ch.OwnedUnitsArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: &ch.UUIDArgs{
					UUID: unit.Player,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		choices := ch.NewChoices(choose1, choose2)
		bases, err := choices.Retrieve(context.Background(), engine, state)
		if err != nil {
			return errors.Wrap(err)
		}
		if len(bases) <= 1 {
			if err := engine.Do(context.Background(), &Event{
				uuid: state.Gen.New(st.EventUUID),
				typ:  EndGameEvent,
				args: &EndGameArgs{
					Winner: state.GetOpponent(unit.Player),
				},
				affect: EndGameAffect,
			}, state); err != nil {
				return errors.Wrap(err)
			}
		}
	}
	return nil
}
