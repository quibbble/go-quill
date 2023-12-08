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

const (
	DrainBaseManaEvent = "DrainBaseMana"
)

type DrainBaseManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func DrainBaseManaAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a DrainBaseManaArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Mana[playerChoice].BaseAmount -= a.Amount
	if state.Mana[playerChoice].BaseAmount < 0 {
		state.Mana[playerChoice].BaseAmount = 0
	}
	return nil
}
