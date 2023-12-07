package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DiscardCardEvent = "DiscardCard"
)

type DiscardCardArgs struct {
	ChoosePlayer parse.Choose
	ChooseCard   parse.Choose
}

func DiscardCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DiscardCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

	cardChoice, err := GetChoice(engine, state, a.ChooseCard, targets...)
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
