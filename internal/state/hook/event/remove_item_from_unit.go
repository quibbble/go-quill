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
	RemoveItemFromUnitEvent = "RemoveItemFromUnit"
)

type RemoveItemFromUnitArgs struct {
	ChooseItem Choose
	ChooseUnit Choose
}

func RemoveItemFromUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a RemoveItemFromUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	chooseItem, err := ch.NewChoice(a.ChooseItem.Type, a.ChooseItem.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	chooseUnit, err := ch.NewChoice(a.ChooseUnit.Type, a.ChooseUnit.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	itemChoices, err := chooseItem.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoices, err := chooseUnit.Retrieve(engine, state, targets...)
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
	item, err := state.Board.XYs[x][y].Unit.(*cd.UnitCard).GetAndRemoveItem(engine, itemChoices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	for _, trait := range item.HeldTraits {
		event := &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  RemoveTraitFromCard,
			args: &RemoveTraitFromCardArgs{
				Trait: trait.GetUUID(),
				Choose: Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: unitChoices[0],
					},
				},
			},
			affect: AddTraitToCardAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}

		// if unit died from adding trait then break
		_, _, err := state.Board.GetUnitXY(unitChoices[0])
		if err != nil {
			break
		}
	}
	return nil
}
