package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	SwapUnitsEvent = "SwapUnits"
)

type SwapUnitsArgs struct {
	ChooseUnitA parse.Choose
	ChooseUnitB parse.Choose
}

func SwapUnitsAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*SwapUnitsArgs)
	unitAChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnitA, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitBChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnitB, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	aX, aY, err := state.Board.GetUnitXY(unitAChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	bX, bY, err := state.Board.GetUnitXY(unitBChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unitA := state.Board.XYs[aX][aY].Unit
	unitB := state.Board.XYs[bX][bY].Unit
	state.Board.XYs[aX][aY].Unit = unitB
	state.Board.XYs[bX][bY].Unit = unitA

	// friends/enemies trait check
	FriendsTraitCheck(e, engine, state)
	EnemiesTraitCheck(e, engine, state)

	return nil
}
