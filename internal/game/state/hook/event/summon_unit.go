package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	SummonUnitEvent = "SummonUnit"
)

type SummonUnitArgs struct {
	ChoosePlayer parse.Choose
	ID           string
	ChooseTile   parse.Choose
}

func SummonUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a SummonUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := GetTileChoice(engine, state, a.ChooseTile, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

	unit, err := state.BuildCard(a.ID, playerChoice)
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
	state.Board.XYs[tX][tY].Unit = unit

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
