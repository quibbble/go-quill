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

const UnitMissingCondition = "UnitMissing"

type UnitMissingArgs struct {
	ChooseUnit parse.Choose
}

func PassUnitMissing(engine *en.Engine, state *st.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	var p UnitMissingArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	unitChoice, err := ev.GetUnitChoice(engine, state, p.ChooseUnit, targets...)
	if err != nil {
		return false, errors.Wrap(err)
	}

	_, _, err = state.Board.GetUnitXY(unitChoice)
	return err != nil, nil
}
