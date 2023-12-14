package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	SummonUnitEvent = "SummonUnit"
)

type SummonUnitArgs struct {
	ChoosePlayer parse.Choose
	ID           string
	ChooseTile   parse.Choose
}

func SummonUnitAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*SummonUnitArgs)
	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := ch.GetTileChoice(ctx, a.ChooseTile, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	unit, err := state.BuildCard(a.ID, playerChoice, true)
	if err != nil {
		return errors.Wrap(err)
	}
	tX, tY, err := state.Board.GetTileXY(tileChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[tX][tY].Unit != nil {
		return errors.Errorf("unit '%s' cannot be placed on a full tile", unit.GetUUID())
	}
	min, max := state.Board.GetPlayableRowRange(playerChoice)
	if tY < min || tY > max {
		return errors.Errorf("unit '%s' must be placed within rows %d to %d", unit.GetUUID(), min, max)
	}

	// haste trait check
	if len(unit.GetTraits(tr.HasteTrait)) > 0 {
		unit.(*cd.UnitCard).Cooldown = 0
	}

	state.Board.XYs[tX][tY].Unit = unit

	// friends/enemies trait check
	FriendsTraitCheck(e, engine, state)
	EnemiesTraitCheck(e, engine, state)

	return nil
}
