package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const UnitMissingCondition = "UnitMissing"

type UnitMissingArgs struct {
	ChooseUnit parse.Choose
}

func PassUnitMissing(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	p := args.(*UnitMissingArgs)

	unitChoice, err := ch.GetUnitChoice(ctx, p.ChooseUnit, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	_, _, err = state.Board.GetUnitXY(unitChoice)
	return err != nil, nil
}
