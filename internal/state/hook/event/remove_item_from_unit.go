package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	RemoveItemFromUnitEvent = "remove_item_from_unit"
)

type RemoveItemFromUnitArgs struct {
	ItemChoice ch.Choose
	UnitChoice ch.Choose
}

func RemoveItemFromUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(RemoveItemFromUnitArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	itemChoices, err := a.ItemChoice.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoices, err := a.UnitChoice.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	if len(itemChoices) != 1 || len(unitChoices) != 1 {
		return errors.ErrInvalidSliceLength
	}
	if itemChoices[0].Type() != st.ItemUUID {
		return st.ErrInvalidUUIDType(itemChoices[0], st.ItemUUID)
	}
	if unitChoices[0].Type() != st.UnitUUID {
		return st.ErrInvalidUUIDType(unitChoices[0], st.UnitUUID)
	}
	x, y, err := state.Board.GetUnitXY(unitChoices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	if err := state.Board.XYs[x][y].Unit.RemoveItem(engine, itemChoices[0]); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
