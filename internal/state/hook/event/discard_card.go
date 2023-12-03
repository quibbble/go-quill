package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DiscardCardEvent = "DiscardCard"
)

type DiscardCardArgs struct {
	Player uuid.UUID
	Choose Choose
}

func DiscardCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DiscardCardArgs
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
	card, err := state.Hand[a.Player].GetCard(choices[0])
	if err != nil {
		return errors.Wrap(err)
	}
	if err := state.Hand[a.Player].RemoveCard(choices[0]); err != nil {
		return errors.Wrap(err)
	}
	if item, ok := card.(*cd.ItemCard); ok {
		item.Reset(state.BuildCard)
	} else if spell, ok := card.(*cd.SpellCard); ok {
		spell.Reset(state.BuildCard)
	} else if unit, ok := card.(*cd.UnitCard); ok {
		unit.Reset(state.BuildCard)
	}
	state.Discard[a.Player].Add(card)
	return nil
}
