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
	AddItemToUnitEvent = "AddItemToUnit"
)

type AddItemToUnitArgs struct {
	Player     uuid.UUID
	ChooseItem Choose
	ChooseUnit Choose
}

func AddItemToUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a AddItemToUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	chooseItem, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), a.ChooseItem.Type, a.ChooseItem.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	chooseUnit, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), a.ChooseUnit.Type, a.ChooseUnit.Args)
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
	card, err := state.Hand[a.Player].GetCard(itemChoices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	item := card.(*cd.ItemCard)
	x, y, err := state.Board.GetUnitXY(unitChoices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[x][y].Unit.(*cd.UnitCard).Type == cd.StructureUnit {
		return errors.Errorf("cannot add an item to a structure unit")
	}
	if err := state.Hand[a.Player].RemoveCard(item.UUID); err != nil {
		return errors.Wrap(err)
	}
	if err := state.Board.XYs[x][y].Unit.(*cd.UnitCard).AddItem(engine, item); err != nil {
		return errors.Wrap(err)
	}
	for _, trait := range item.HeldTraits {
		event := &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  AddTraitToCard,
			args: &AddTraitToCardArgs{
				Trait: Trait{
					Type: trait.GetType(),
					Args: trait.GetArgs(),
				},
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
