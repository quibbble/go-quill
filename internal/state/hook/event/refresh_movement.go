package event

import (
	"github.com/mitchellh/mapstructure"
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	RefreshMovementEvent = "RefreshMovement"
)

type RefreshMovementArgs struct {
	Choose Choose
}

func RefreshMovementAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a RefreshMovementArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoice(a.Choose.Type, a.Choose.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		x, y, err := state.Board.GetUnitXY(choice)
		if err != nil {
			return errors.Wrap(err)
		}
		unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
		unit.Movement = unit.GetInit().(*cards.UnitCard).Movement
	}
	return nil
}
