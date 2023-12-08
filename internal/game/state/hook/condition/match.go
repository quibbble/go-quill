package condition

import (
	"context"

	"github.com/mitchellh/mapstructure"
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
	var p MatchArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}
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
