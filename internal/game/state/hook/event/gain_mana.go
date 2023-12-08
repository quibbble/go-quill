package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	GainManaEvent = "GainMana"
)

type GainManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func GainManaAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a GainManaArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Mana[playerChoice].Amount += a.Amount
	return nil
}
