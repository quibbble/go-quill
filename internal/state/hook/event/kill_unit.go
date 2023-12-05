package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	KillUnitEvent = "KillUnit"
)

type KillUnitArgs struct {
	Choose Choose
}

func KillUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a KillUnitArgs
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
	state.Board.XYs[x][y] = nil

	// death cry trait check
	deathCrys := unit.GetTraits(tr.DeathCryTrait)
	for _, deathCry := range deathCrys {
		args := deathCry.GetArgs().(*tr.DeathCryArgs)
		event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	if unit.GetInit().(*parse.UnitCard).ID != "U0001" {
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
		choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), ch.BasesChoice, ch.BasesArgs{
			Players: []uuid.UUID{unit.Player},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		bases, err := choose.Retrieve(engine, state)
		if err != nil {
			return errors.Wrap(err)
		}
		if len(bases) <= 1 {
			if err := engine.Do(&Event{
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
