package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	PlaceUnitEvent = "PlaceUnit"
)

type PlaceUnitArgs struct {
	X, Y   int
	Player uuid.UUID
	Choose Choose
}

func PlaceUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a PlaceUnitArgs
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
	if len(choices) != 1 {
		return errors.ErrInvalidSliceLength
	}
	if choices[0].Type() != st.UnitUUID {
		return st.ErrInvalidUUIDType(choices[0], st.UnitUUID)
	}
	card, err := state.Hand[a.Player].GetCard(choices[0])
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
	if err := state.Hand[a.Player].RemoveCard(choices[0]); err != nil {
		return errors.Wrap(err)
	}
	state.Board.XYs[a.X][a.Y].Unit = unit
	return nil
}
