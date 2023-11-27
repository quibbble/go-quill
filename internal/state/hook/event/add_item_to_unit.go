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
	AddItemToUnitEvent = "AddItemToUnit"
)

type AddItemToUnitArgs struct {
	Player     uuid.UUID
	ItemChoice ch.Choose
	UnitChoice ch.Choose
}

func AddItemToUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(AddItemToUnitArgs)
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
	card, err := state.Hand[a.Player].GetCard(itemChoices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	item := card.(*cd.ItemCard)
	x, y, err := state.Board.GetUnitXY(unitChoices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	if err := state.Hand[a.Player].RemoveCard(item.UUID); err != nil {
		return errors.Wrap(err)
	}
	if err := state.Board.XYs[x][y].Unit.AddItem(engine, item); err != nil {
		return errors.Wrap(err)
	}
	return nil
}
