package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	DiscardCardEvent = "DiscardCard"
)

type DiscardCardArgs struct {
	ChoosePlayer parse.Choose
	ChooseCard   parse.Choose
}

func DiscardCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	a := args.(*DiscardCardArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	cardChoice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(cardChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	if err := state.Hand[playerChoice].RemoveCard(cardChoice); err != nil {
		return errors.Wrap(err)
	}
	if item, ok := card.(*cd.ItemCard); ok {
		item.Reset(state.BuildCard)
	} else if spell, ok := card.(*cd.SpellCard); ok {
		spell.Reset(state.BuildCard)
	} else if unit, ok := card.(*cd.UnitCard); ok {
		unit.Reset(state.BuildCard)
	}
	state.Discard[playerChoice].Add(card)
	return nil
}
