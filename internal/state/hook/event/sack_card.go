package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	SackCardEvent = "sack_card"

	ManaSackOption  = "mana"
	CardsSackOption = "cards"
)

type SackCardArgs struct {
	Player uuid.UUID
	ch.Choose
	Option string
}

func SackCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(SackCardArgs)
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

	events := []*Event{
		{
			uuid: uuid.New(st.EventUUID),
			typ:  DiscardCardEvent,
			args: &DiscardCardArgs{
				Player: a.Player,
				Choose: a.Choose,
			},
			affect: DiscardCardAffect,
		},
	}

	switch a.Option {
	case ManaSackOption:
		events = append(events, &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  GainBaseManaEvent,
			args: &GainBaseManaArgs{
				Player: a.Player,
				Amount: 1,
			},
			affect: GainBaseManaAffect,
		}, &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  GainManaEvent,
			args: &GainManaArgs{
				Player: a.Player,
				Amount: 1,
			},
			affect: GainManaAffect,
		})
	case CardsSackOption:
		events = append(events, &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				Player: a.Player,
			},
			affect: DrawCardAffect,
		}, &Event{
			uuid: uuid.New(st.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				Player: a.Player,
			},
			affect: DrawCardAffect,
		})
	default:
		return errors.Errorf("invalid sack option '%s'", a.Option)
	}

	for _, event := range events {
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
