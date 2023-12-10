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
}

func RecallUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	a := args.(*RecallUnitArgs)
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

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	// reset and move items and unit back to hand
	for _, item := range unit.Items {
		item.Reset(state.BuildCard)
		state.Discard[item.Player].Add(item)
	}
	unit.Reset(state.BuildCard)
	state.Hand[unit.Player].Add(unit)
	return nil
}
