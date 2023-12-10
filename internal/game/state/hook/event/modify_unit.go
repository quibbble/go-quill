package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
)

const (
	ModifyUnitEvent = "ModifyUnit"
)

type ModifyUnitArgs struct {
	ChooseUnit parse.Choose
	Stat       string
	Amount     int
}

func ModifyUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	a := args.(*ModifyUnitArgs)
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	x, y, err := state.Board.GetUnitXY(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

	switch a.Stat {
	case cd.MovementStat:
		unit.Movement = maths.MaxInt(0, unit.Movement+a.Amount)
	case cd.CooldownStat:
		unit.Cooldown = maths.MaxInt(0, unit.Cooldown+a.Amount)
	default:
		return errors.Errorf("'%s' is not a stat that may be modified", a.Stat)
	}

	return nil
}
