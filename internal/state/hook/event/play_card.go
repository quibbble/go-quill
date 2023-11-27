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
	PlayCardEvent = "PlayCard"
)

type PlayCardArgs struct {
	Player uuid.UUID
	ch.Choose
}

func PlayCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(PlayCardArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	choices, err := a.Choose.Retrieve(engine, state, targets...)
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

	if err := engine.Do(&Event{
		uuid: uuid.New(st.EventUUID),
		typ:  DrainManaEvent,
		args: &DrainManaArgs{
			Player: a.Player,
			Amount: card.Cost,
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
				ItemChoice: &ch.UUIDChoice{
					UUID: card.UUID,
				},
				UnitChoice: &ch.UUIDChoice{
					UUID: targets[0],
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
				Choose: &ch.UUIDChoice{
					UUID: card.UUID,
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
				Choose: &ch.UUIDChoice{
					UUID: card.UUID,
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
