package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	MoveUnitEvent = "MoveUnit"
)

type MoveUnitArgs struct {
	ChooseUnit parse.Choose
	ChooseTile parse.Choose
}

func MoveUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a MoveUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	unitChoice, err := GetUnitChoice(engine, state, a.ChooseUnit, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := GetTileChoice(engine, state, a.ChooseTile, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

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
	if !unit.CheckCodex(uX, uY, tX, tY) {
		return errors.Errorf("unit '%s' cannot move due to failed codex check", unit.UUID)
	}
	if unit.Movement < 1 {
		return errors.Errorf("unit '%s' cannot move with no movement", unit.UUID)
	}
	state.Board.XYs[uX][uY].Unit = nil
	state.Board.XYs[tX][tY].Unit = unit
	unit.Movement--

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
