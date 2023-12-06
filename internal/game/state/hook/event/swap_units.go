package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	SwapUnitsEvent = "SwapUnits"
)

type SwapUnitsArgs struct {
	UnitA uuid.UUID
	UnitB uuid.UUID
}

func SwapUnitsAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a SwapUnitsArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	aX, aY, err := state.Board.GetUnitXY(a.UnitA)
	if err != nil {
		return errors.Wrap(err)
	}
	bX, bY, err := state.Board.GetUnitXY(a.UnitA)
	if err != nil {
		return errors.Wrap(err)
	}
	unitA := state.Board.XYs[aX][aY].Unit
	unitB := state.Board.XYs[bX][bY].Unit
	state.Board.XYs[aX][aY].Unit = unitB
	state.Board.XYs[bX][bY].Unit = unitA

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
