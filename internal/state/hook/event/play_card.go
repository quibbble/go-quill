package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	PlayCardEvent = "PlayCard"
)

type PlayCardArgs struct {
	Player uuid.UUID
	Choose Choose
}

func PlayCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a PlayCardArgs
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
	i, err := state.Hand[a.Player].GetCard(choices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	card := i.(*cd.Card)

	// drain mana equal to card cost
	if state.Mana[a.Player].Amount < card.Cost {
		return errors.Errorf("player '%s' does not have enough mana to play '%s'", a.Player, card.UUID)
	}

	// purity trait check
	if card.UUID.Type() == st.SpellUUID {
		for _, target := range targets {
			if x, y, err := state.Board.GetUnitXY(target); err == nil {
				unit := state.Board.XYs[x][y].Unit
				if len(unit.GetTraits(tr.PurityTrait)) > 0 {
					return errors.Errorf("'%s' cannot target '%s' due to purity trait", card.UUID, unit.GetUUID())
				}
			}

		}
	}

	if err := engine.Do(&Event{
		uuid: uuid.New(st.EventUUID),
		typ:  DrainManaEvent,
		args: &DrainManaArgs{
			Player: a.Player,
			Amount: maths.MaxInt(card.Cost, 0),
		},
		affect: DrainManaAffect,
	}, state); err != nil {
		return errors.Wrap(err)
	}

	// add hooks
	for _, hook := range card.Hooks {
		engine.Register(hook)
	}

	// create event based on card type
	var event *Event
	switch card.UUID.Type() {
	case st.ItemUUID:
		if len(targets) <= 0 {
			return errors.ErrIndexOutOfBounds
		}
		event = &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  AddItemToUnitEvent,
			args: &AddItemToUnitArgs{
				Player: a.Player,
				ChooseItem: Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.UUID,
					},
				},
				ChooseUnit: Choose{
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
			uuid: uuid.New(st.EventUUID),
			typ:  DiscardCardEvent,
			args: &DiscardCardArgs{
				Player: a.Player,
				Choose: Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.UUID,
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
			uuid: uuid.New(st.EventUUID),
			typ:  PlaceUnitEvent,
			args: &PlaceUnitArgs{
				X:      x,
				Y:      y,
				Player: a.Player,
				Choose: Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: card.UUID,
					},
				},
			},
			affect: PlaceUnitAffect,
		}
	}

	// apply card event and then any additional events attached to the card
	events := append([]en.IEvent{event}, card.Events...)
	for _, event := range events {
		if err := engine.Do(event, state, targets...); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
