package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	HealUnitEvent = "HealUnit"
)

type HealUnitArgs struct {
	Amount     int
	ChooseUnit parse.Choose
}

func HealUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a HealUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
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
