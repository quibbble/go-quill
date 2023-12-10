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
	GainBaseManaEvent = "GainBaseMana"
)

type GainBaseManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func GainBaseManaAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	a := args.(*GainBaseManaArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	state.Mana[playerChoice].BaseAmount += a.Amount
	return nil
}
