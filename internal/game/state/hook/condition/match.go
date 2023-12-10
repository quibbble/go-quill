package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const MatchCondition = "Match"

type MatchArgs struct {
	ChooseA parse.Choose
	ChooseB parse.Choose
}

func PassMatch(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	p := args.(*MatchArgs)
	a, err := choose.GetChoice(ctx, p.ChooseA, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	b, err := choose.GetChoice(ctx, p.ChooseB, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	return a == b, nil
}
