package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const SackCardEvent = "SackCard"

const (
	ManaSackOption  = "Mana"
	CardsSackOption = "Cards"
)

type SackCardArgs struct {
	Player uuid.UUID
	Option string
	Choose Choose
}

func SackCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a SackCardArgs
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

	events := []*Event{
		{
			uuid: state.Gen.New(st.EventUUID),
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
			uuid: state.Gen.New(st.EventUUID),
			typ:  GainBaseManaEvent,
			args: &GainBaseManaArgs{
				Player: a.Player,
				Amount: 1,
			},
			affect: GainBaseManaAffect,
		}, &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  GainManaEvent,
			args: &GainManaArgs{
				Player: a.Player,
				Amount: 1,
			},
			affect: GainManaAffect,
		})
	case CardsSackOption:
		events = append(events, &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				Player: a.Player,
			},
			affect: DrawCardAffect,
		}, &Event{
			uuid: state.Gen.New(st.EventUUID),
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
