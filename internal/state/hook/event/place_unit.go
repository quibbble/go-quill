package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	PlaceUnitEvent = "place_unit"
)

type PlaceUnitArgs struct {
	X, Y   int
	Player uuid.UUID
	ch.Choose
}

func PlaceUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(PlaceUnitArgs)
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
	card, err := state.Hand[a.Player].GetAndRemoveCard(choices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	unit := card.(*cd.UnitCard)
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
