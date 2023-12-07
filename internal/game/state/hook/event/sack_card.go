package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const SackCardEvent = "SackCard"

const (
	ManaSackOption  = "Mana"
	CardsSackOption = "Cards"
)

type SackCardArgs struct {
	ChoosePlayer parse.Choose
	Option       string
	ChooseCard   parse.Choose
}

func SackCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a SackCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

	events := []*Event{
		{
			uuid: state.Gen.New(st.EventUUID),
			typ:  DiscardCardEvent,
			args: &DiscardCardArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.CurrentPlayerChoice,
					Args: &ch.CurrentPlayerArgs{},
				},
				ChooseCard: a.ChooseCard,
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
				ChoosePlayer: parse.Choose{
					Type: ch.CurrentPlayerChoice,
					Args: &ch.CurrentPlayerArgs{},
				},
				Amount: 1,
			},
			affect: GainBaseManaAffect,
		}, &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  GainManaEvent,
			args: &GainManaArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.CurrentPlayerChoice,
					Args: &ch.CurrentPlayerArgs{},
				},
				Amount: 1,
			},
			affect: GainManaAffect,
		})
	case CardsSackOption:
		events = append(events, &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: playerChoice,
					},
				},
			},
			affect: DrawCardAffect,
		}, &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  DrawCardEvent,
			args: &DrawCardArgs{
				ChoosePlayer: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: playerChoice,
					},
				},
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

	state.Sacked[playerChoice] = true

	return nil
}
