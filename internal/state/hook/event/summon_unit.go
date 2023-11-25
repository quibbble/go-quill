package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	SummonUnitEvent = "summon_unit"
)

type SummonUnitArgs struct {
	X, Y   int
	Player uuid.UUID
	ID     string
}

func SummonUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(SummonUnitArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	unit, err := card.NewUnitCard(a.ID, a.Player)
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[a.X][a.Y].Unit != nil {
		return errors.Errorf("unit '%s' cannot be placed on a full tile", unit.UUID)
	}
	min, max := state.Board.GetPlayableRowRange(a.Player)
	if a.Y < min || a.Y > max {
		return errors.Errorf("unit '%s' must be placed within rows %d to %d", unit.UUID, min, max)
	}
	state.Board.XYs[a.X][a.Y].Unit = unit
	return nil
}
