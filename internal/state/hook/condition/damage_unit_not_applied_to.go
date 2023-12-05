package condition

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	ev "github.com/quibbble/go-quill/internal/state/hook/event"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const DamageUnitAppliedToCondition = "DamageUnitAppliedTo"

type DamageUnitAppliedToArgs struct {
	Choose ch.RawChoose
}

func PassDamageUnitAppliedTo(engine *en.Engine, state *st.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	var p DamageUnitAppliedToArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), p.Choose.Type, p.Choose.Type)
	if err != nil {
		return false, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return false, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return false, errors.ErrInvalidSliceLength
	}
	conditionUnit := choices[0]

	if event == nil {
		return false, errors.ErrNilInterface
	}
	a, ok := event.GetArgs().(*ev.DamageUnitArgs)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	choose, err = ch.NewChoose(state.Gen.New(st.ChooseUUID), a.Choose.Type, a.Choose.Type)
	if err != nil {
		return false, errors.Wrap(err)
	}
	choices, err = choose.Retrieve(engine, state, targets...)
	if err != nil {
		return false, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return false, errors.ErrInvalidSliceLength
	}
	eventUnit := choices[0]
	return conditionUnit == eventUnit, nil
}
