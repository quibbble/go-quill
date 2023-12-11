package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const ManaAboveCondition = "ManaAbove"

type ManaAboveArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func PassManaAbove(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	p := c.GetArgs().(*ManaAboveArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, p.ChoosePlayer, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	return state.Mana[playerChoice].Amount > p.Amount, nil
}
