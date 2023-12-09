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

const DrawCardEvent = "DrawCard"

type DrawCardArgs struct {
	ChoosePlayer parse.Choose
}

func DrawCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a DrawCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	if state.Hand[playerChoice].GetSize() <= st.MaxHandSize {
		card, err := state.Deck[playerChoice].Draw()
		if err != nil {
			return errors.Wrap(err)
		}
		state.Hand[playerChoice].Add(*card)
	} else {
		event := &Event{
			uuid: state.Gen.New(en.EventUUID),
			typ:  BurnCardEvent,
			args: BurnCardArgs{
				ChoosePlayer: a.ChoosePlayer,
			},
			affect: BurnCardAffect,
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
