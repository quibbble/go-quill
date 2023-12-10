package condition

import (
	"context"
	"slices"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const ContainsCondition = "Contains"

type ContainsArgs struct {
	ChooseChain parse.Choose
	Choose      parse.Choose
}

func PassContains(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	p := args.(*ContainsArgs)
	choices, err := ch.GetChoices(ctx, p.ChooseChain, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	choice, err := ch.GetChoice(ctx, p.Choose, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	return slices.Contains(choices, choice), nil
}
