package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	AddItemToUnitEvent = "AddItemToUnit"
)

type AddItemToUnitArgs struct {
	ChoosePlayer parse.Choose
	ChooseItem   parse.Choose
	ChooseUnit   parse.Choose
}

func AddItemToUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a AddItemToUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	itemChoice, err := GetItemChoice(engine, state, a.ChooseItem, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := GetUnitChoice(engine, state, a.ChooseUnit, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(itemChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	itemCard := card.(*cd.ItemCard)

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unitCard := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

	if unitCard.Type == cd.StructureUnit {
		return errors.Errorf("cannot add an item to a structure unit")
	}
	if err := state.Hand[playerChoice].RemoveCard(itemChoice); err != nil {
		return errors.Wrap(err)
	}
	if err := unitCard.AddItem(engine, itemCard); err != nil {
		return errors.Wrap(err)
	}

	for _, trait := range itemCard.HeldTraits {
		event := &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  AddTraitToCard,
			args: &AddTraitToCardArgs{
				Trait: parse.Trait{
					Type: trait.GetType(),
					Args: trait.GetArgs(),
				},
				ChooseCard: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: unitChoice,
					},
				},
			},
			affect: AddTraitToCardAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}

		// if unit died from adding trait then break
		_, _, err := state.Board.GetUnitXY(unitChoice)
		if err != nil {
			break
		}
	}
	return nil
}
