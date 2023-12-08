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

const DamageUnitAppliedToCondition = "DamageUnitAppliedTo"

type DamageUnitAppliedToArgs struct {
	ChooseUnit parse.Choose
}

func PassDamageUnitAppliedTo(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	var p DamageUnitAppliedToArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	conditionUnit, err := ev.GetUnitChoice(ctx, p.ChooseUnit, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	event := ctx.Value(en.HookEventCtx).(en.IEvent)

	if event == nil {
		return false, errors.ErrNilInterface
	}

	var a ev.DamageUnitArgs
	if err := mapstructure.Decode(event.GetArgs(), &a); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	eventUnit, err := ev.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	return conditionUnit == eventUnit, nil
}
