package condition

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
)

const UnitMissingCondition = "UnitMissing"

type UnitMissingArgs struct {
	Choose ch.RawChoose
}

func PassUnitMissing(engine *en.Engine, state *st.State, args interface{}, event ...en.IEvent) (bool, error) {
	var p UnitMissingArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), p.Choose.Type, p.Choose.Type)
	if err != nil {
		return false, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return false, errors.ErrInvalidSliceLength
	}
	if choices[0].Type() != st.UnitUUID {
		return false, st.ErrInvalidUUIDType(choices[0], st.UnitUUID)
	}
	_, _, err = state.Board.GetUnitXY(choices[0])
	return err != nil, nil
}
