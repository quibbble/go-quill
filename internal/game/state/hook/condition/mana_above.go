package condition

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ev "github.com/quibbble/go-quill/internal/game/state/hook/event"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const ManaAboveCondition = "ManaAbove"

type ManaAboveArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func PassManaAbove(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	var p ManaAboveArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	playerChoice, err := ev.GetPlayerChoice(ctx, p.ChoosePlayer, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	return state.Mana[playerChoice].Amount > p.Amount, nil
}
