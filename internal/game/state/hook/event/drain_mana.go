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
	DrainManaEvent = "DrainMana"
)

type DrainManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func DrainManaAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*DrainManaArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Mana[playerChoice].Amount -= a.Amount
	if state.Mana[playerChoice].Amount < 0 {
		state.Mana[playerChoice].Amount = 0
	}
	return nil
}
