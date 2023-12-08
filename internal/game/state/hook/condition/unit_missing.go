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

const UnitMissingCondition = "UnitMissing"

type UnitMissingArgs struct {
	ChooseUnit parse.Choose
}

func PassUnitMissing(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	var p UnitMissingArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	unitChoice, err := ev.GetUnitChoice(ctx, p.ChooseUnit, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	_, _, err = state.Board.GetUnitXY(unitChoice)
	return err != nil, nil
}
