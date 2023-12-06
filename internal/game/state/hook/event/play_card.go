package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	PlayCardEvent = "PlayCard"
)

type PlayCardArgs struct {
	Player uuid.UUID
	Choose ch.RawChoose
}

func PlayCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a PlayCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), a.Choose.Type, a.Choose.Args)
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
	card, err := state.Hand[a.Player].GetCard(choices[0])
	if err != nil {
		return errors.Wrap(err)
	}

	// drain mana equal to card cost
	if state.Mana[a.Player].Amount < card.GetCost() {
		return errors.Errorf("player '%s' does not have enough mana to play '%s'", a.Player, card.GetUUID())
	}

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

	if err := engine.Do(&Event{
		uuid: state.Gen.New(st.EventUUID),
		typ:  DrainManaEvent,
		args: &DrainManaArgs{
			Player: a.Player,
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
				Player: a.Player,
				ChooseItem: ch.RawChoose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.GetUUID(),
					},
				},
				ChooseUnit: ch.RawChoose{
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
				Player: a.Player,
				Choose: ch.RawChoose{
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
		x, y, err := state.Board.GetTileXY(targets[0])
		if err != nil {
			return errors.Wrap(err)
		}
		event = &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  PlaceUnitEvent,
			args: &PlaceUnitArgs{
				X:      x,
				Y:      y,
				Player: a.Player,
				Choose: ch.RawChoose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.GetUUID(),
					},
				},
			},
			affect: PlaceUnitAffect,
		}
	}

	// apply card event and then any additional events attached to the card
	events := append([]en.IEvent{event}, card.GetEvents()...)
	for _, event := range events {
		if err := engine.Do(event, state, targets...); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
