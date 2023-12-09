package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	MoveUnitEvent = "MoveUnit"
)

type MoveUnitArgs struct {
	ChooseUnit   parse.Choose
	ChooseTile   parse.Choose
	UnitMovement bool
}

func MoveUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a MoveUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	choose, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), a.ChooseTile.Type, a.ChooseTile.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	if len(choices) == 0 {
		return nil
	}
	tileChoice := choices[0]

	uX, uY, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[uX][uY].Unit.(*cd.UnitCard)

	tX, tY, err := state.Board.GetTileXY(tileChoice)
	if err != nil {
		return errors.Wrap(err)
	}

	if state.Board.XYs[tX][tY].Unit != nil {
		return errors.Errorf("unit '%s' cannot move to a full tile", unit.UUID)
	}

	if a.UnitMovement {
		if !unit.CheckCodex(uX, uY, tX, tY) {
			return errors.Errorf("unit '%s' cannot move due to failed codex check", unit.UUID)
		}
		if unit.Movement < 1 {
			return errors.Errorf("unit '%s' cannot move with no movement", unit.UUID)
		}
		unit.Movement--
	}

	state.Board.XYs[uX][uY].Unit = nil
	state.Board.XYs[tX][tY].Unit = unit

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
