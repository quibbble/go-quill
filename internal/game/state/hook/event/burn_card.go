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
	BurnCardEvent = "BurnCard"
)

type BurnCardArgs struct {
	ChoosePlayer parse.Choose
}

func BurnCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a BurnCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Deck[playerChoice].Draw()
	if err != nil {
		return errors.Wrap(err)
	}
	state.Trash[playerChoice].Add(*card)
	return nil
}
