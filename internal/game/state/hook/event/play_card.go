package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	PlayCardEvent = "PlayCard"
)

type PlayCardArgs struct {
	ChoosePlayer parse.Choose
	ChooseCard   parse.Choose
}

func PlayCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a PlayCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	cardChoice, err := GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(cardChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	playable, err := card.Playable(engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	next, err := card.NextTargets(ctx, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	if !playable || len(next) != 0 {
		return errors.Errorf("card cannot be played")
	}

	// drain mana equal to card cost
	if state.Mana[playerChoice].Amount < card.GetCost() {
		return errors.Errorf("player '%s' does not have enough mana to play '%s'", playerChoice, card.GetUUID())
	}

	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)

	// purity trait check
	if card.GetUUID().Type() == st.SpellUUID {
		for _, target := range targets {
			if x, y, err := state.Board.GetUnitXY(target); err == nil {
				unit := state.Board.XYs[x][y].Unit
				if len(unit.GetTraits(tr.PurityTrait)) > 0 {
					return errors.Errorf("'%s' cannot target '%s' due to purity trait", card.GetUUID(), unit.GetUUID())
				}
			}

		}
	}

	if err := engine.Do(context.Background(), &Event{
		uuid: state.Gen.New(st.EventUUID),
		typ:  DrainManaEvent,
		args: &DrainManaArgs{
			ChoosePlayer: parse.Choose{
				Type: ch.CurrentPlayerChoice,
				Args: &ch.CurrentPlayerArgs{},
			},
			Amount: maths.MaxInt(card.GetCost(), 0),
		},
		affect: DrainManaAffect,
	}, state); err != nil {
		return errors.Wrap(err)
	}

	// add hooks
	for _, hook := range card.GetHooks() {
		engine.Register(hook)
	}

	// create event based on card type
	var event *Event
	switch card.GetUUID().Type() {
	case st.ItemUUID:
		if len(targets) <= 0 {
			return errors.ErrIndexOutOfBounds
		}
		event = &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  AddItemToUnitEvent,
			args: &AddItemToUnitArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: playerChoice,
					},
				},
				ChooseItem: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.GetUUID(),
					},
				},
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: targets[0],
					},
				},
			},
			affect: AddItemToUnitAffect,
		}
	case st.SpellUUID:
		event = &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  DiscardCardEvent,
			args: &DiscardCardArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.CurrentPlayerChoice,
					Args: &ch.CurrentPlayerArgs{},
				},
				ChooseCard: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.GetUUID(),
					},
				},
			},
			affect: DiscardCardAffect,
		}
	case st.UnitUUID:
		if len(targets) <= 0 {
			return errors.ErrIndexOutOfBounds
		}
		event = &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  PlaceUnitEvent,
			args: &PlaceUnitArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.CurrentPlayerChoice,
					Args: &ch.CurrentPlayerArgs{},
				},
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.GetUUID(),
					},
				},
				ChooseTile: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: targets[0],
					},
				},
			},
			affect: PlaceUnitAffect,
		}
	}

	// apply card event and then any additional events attached to the card
	events := append([]en.IEvent{event}, card.GetEvents()...)
	for _, event := range events {
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
