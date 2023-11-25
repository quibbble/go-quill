package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	KillUnitEvent = "kill_unit"
)

type KillUnitArgs struct {
	ch.Choose
}

func KillUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(KillUnitArgs)
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
	state.Board.XYs[x][y] = nil
	state.Discard[unit.Owner].Add(unit)
	return nil
}
