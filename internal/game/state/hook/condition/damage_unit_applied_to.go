package condition

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ev "github.com/quibbble/go-quill/internal/game/state/hook/event"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const DamageUnitAppliedToCondition = "DamageUnitAppliedTo"

type DamageUnitAppliedToArgs struct {
	ChooseUnit parse.Choose
}

func PassDamageUnitAppliedTo(engine *en.Engine, state *st.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	var p DamageUnitAppliedToArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	conditionUnit, err := ev.GetUnitChoice(engine, state, p.ChooseUnit, targets...)
	if err != nil {
		return false, errors.Wrap(err)
	}

	if event == nil {
		return false, errors.ErrNilInterface
	}

	var a ev.DamageUnitArgs
	if err := mapstructure.Decode(event.GetArgs(), &a); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	eventUnit, err := ev.GetUnitChoice(engine, state, a.ChooseUnit, targets...)
	if err != nil {
		return false, errors.Wrap(err)
	}

	return conditionUnit == eventUnit, nil
}
