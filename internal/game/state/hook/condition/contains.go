package condition

import (
	"context"
	"slices"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const ContainsCondition = "Contains"

type ContainsArgs struct {
	Choices parse.Choose
	Choice  parse.Choose
}

func PassContains(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	var p ContainsArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}
	choices, err := choose.GetChoices(ctx, p.Choices, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	choice, err := choose.GetChoice(ctx, p.Choice, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	return slices.Contains(choices, choice), nil
}
