package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const SackCardEvent = "SackCard"

const (
	ManaSackOption  = "Mana"
	CardsSackOption = "Cards"
)

type SackCardArgs struct {
	ChoosePlayer parse.Choose
	SackOption   string
	ChooseCard   parse.Choose
}

func SackCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a SackCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	events := []*Event{
		{
			uuid: state.Gen.New(en.EventUUID),
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

	switch a.SackOption {
	case ManaSackOption:
		events = append(events, &Event{
			uuid: state.Gen.New(en.EventUUID),
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
			uuid: state.Gen.New(en.EventUUID),
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
			uuid: state.Gen.New(en.EventUUID),
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
			uuid: state.Gen.New(en.EventUUID),
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
		return errors.Errorf("invalid sack option '%s'", a.SackOption)
	}

	for _, event := range events {
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	state.Sacked[playerChoice] = true

	return nil
}
