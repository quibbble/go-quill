package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	RecycleDeckEvent = "RecycleDeck"
)

type RecycleDeckArgs struct {
	ChoosePlayer parse.Choose
}

func RecycleDeckAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RecycleDeckArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Deck[playerChoice] = state.Discard[playerChoice]
	state.Discard[playerChoice] = st.NewEmptyDeck(state.Seed)
	state.Deck[playerChoice].Shuffle()
	return nil
}
