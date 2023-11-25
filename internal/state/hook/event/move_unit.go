package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	MoveUnitEvent = "move_unit"
)

type MoveUnitArgs struct {
	X, Y int
	ch.Choose
}

func MoveUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(MoveUnitArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	choices, err := a.Choose.Retrieve(engine, state, targets...)
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
	unit := state.Board.XYs[x][y].Unit
	if state.Board.XYs[a.X][a.Y].Unit != nil {
		return errors.Errorf("unit '%s' cannot move to a full tile", unit.UUID)
	}
	if !unit.CheckCodex(x, y, a.X, a.Y) {
		return errors.Errorf("unit '%s' cannot move due to failed codex check", unit.UUID)
	}
	if unit.Movement < 1 {
		return errors.Errorf("unit '%s' cannot move with no movement", unit.UUID)
	}
	state.Board.XYs[x][y].Unit = nil
	state.Board.XYs[a.X][a.Y].Unit = unit
	unit.Movement--
	return nil
}
