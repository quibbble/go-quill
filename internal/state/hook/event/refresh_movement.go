package event

import (
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	RefreshMovementEvent = "RefreshMovement"
)

type RefreshMovementArgs struct {
	ch.Choose
}

func RefreshMovementAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(RefreshMovementArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	choices, err := a.Choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		x, y, err := state.Board.GetUnitXY(choice)
		if err != nil {
			return errors.Wrap(err)
		}
		unit := state.Board.XYs[x][y].Unit
		unit.Movement = unit.GetInit().(*cards.UnitCard).Movement
	}
	return nil
}
