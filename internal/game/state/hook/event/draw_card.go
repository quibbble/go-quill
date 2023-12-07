package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const DrawCardEvent = "DrawCard"

type DrawCardArgs struct {
	ChoosePlayer parse.Choose
}

func DrawCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DrawCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
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
			uuid: state.Gen.New(st.EventUUID),
			typ:  BurnCardEvent,
			args: &BurnCardArgs{
				ChoosePlayer: a.ChoosePlayer,
			},
			affect: BurnCardAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
