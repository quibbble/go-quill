package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	HealUnitEvent = "HealUnit"
)

type HealUnitArgs struct {
	Amount     int
	ChooseUnit parse.Choose
}

func HealUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a HealUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	unitChoice, err := GetUnitChoice(engine, state, a.ChooseUnit, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)
	unit.Health += a.Amount
	if unit.Health > unit.GetInit().(*parse.UnitCard).Health {
		unit.Health = unit.GetInit().(*parse.UnitCard).Health
	}
	return nil
}
