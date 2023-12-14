package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	AddTraitToCard = "AddTraitToCard"
)

type AddTraitToCardArgs struct {
	Trait      parse.Trait
	ChooseCard parse.Choose
}

func AddTraitToCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*AddTraitToCardArgs)
	trait, err := tr.NewTrait(state.Gen.New(en.ChooseUUID), a.Trait.Type, a.Trait.Args)
	if err != nil {
		return errors.Wrap(err)
	}

	choice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return st.ErrNotFound(choice)
	}
	if card.AddTrait(trait); err != nil {
		return errors.Wrap(err)
	}

	// friends/enemies trait check
	FriendsTraitCheck(e, engine, state)
	EnemiesTraitCheck(e, engine, state)

	return nil
}
