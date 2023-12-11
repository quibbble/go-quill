package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	RecallUnitEvent = "RecallUnit"
)

type RecallUnitArgs struct {
	ChooseUnit parse.Choose

	// DO NOT SET IN YAML - SET BY ENGINE
	// tile location unit before recall
	ChooseTile parse.Choose
}

func RecallUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RecallUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

	if unit.GetID() == "U0001" {
		return errors.Errorf("cannot rescind U0001")
	}

	state.Board.XYs[x][y].Unit = nil
	a.ChooseTile = parse.Choose{
		Type: ch.UUIDChoice,
		Args: ch.UUIDArgs{
			UUID: state.Board.XYs[x][y].UUID,
		},
	}

	// friends/enemies trait check
	FriendsTraitCheck(e, engine, state)
	EnemiesTraitCheck(e, engine, state)

	// reset and move items and unit back to hand
	for _, item := range unit.Items {
		item.Reset(state.BuildCard)
		state.Discard[item.Player].Add(item)
	}
	unit.Reset(state.BuildCard)

	if state.Hand[unit.Player].GetSize() > st.MaxHandSize {
		// burn the card if hand size to large
		state.Deck[unit.Player].Add(unit)
		event, err := NewEvent(state.Gen.New(en.EventUUID), BurnCardEvent, BurnCardArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.UUIDChoice,
				Args: ch.UUIDArgs{
					UUID: unit.Player,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	} else {
		state.Hand[unit.Player].Add(unit)
	}
	return nil
}
