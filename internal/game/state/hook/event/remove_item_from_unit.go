package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	RemoveItemFromUnitEvent = "RemoveItemFromUnit"
)

type RemoveItemFromUnitArgs struct {
	ChooseItem parse.Choose
	ChooseUnit parse.Choose
}

func RemoveItemFromUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a RemoveItemFromUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	itemChoice, err := ch.GetItemChoice(ctx, a.ChooseItem, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	item, err := state.Board.XYs[x][y].Unit.(*cd.UnitCard).GetAndRemoveItem(engine, itemChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, trait := range item.HeldTraits {
		event := &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  RemoveTraitFromCard,
			args: &RemoveTraitFromCardArgs{
				Trait: trait.GetUUID(),
				ChooseCard: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: unitChoice,
					},
				},
			},
			affect: AddTraitToCardAffect,
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
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
