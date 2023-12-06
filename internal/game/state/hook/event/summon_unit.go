package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	SummonUnitEvent = "SummonUnit"
)

type SummonUnitArgs struct {
	X, Y   int
	Player uuid.UUID
	ID     string
}

func SummonUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a SummonUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	unit, err := state.BuildCard(a.ID, a.Player)
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[a.X][a.Y].Unit != nil {
		return errors.Errorf("unit '%s' cannot be placed on a full tile", unit.GetUUID())
	}
	min, max := state.Board.GetPlayableRowRange(a.Player)
	if a.Y < min || a.Y > max {
		return errors.Errorf("unit '%s' must be placed within rows %d to %d", unit.GetUUID(), min, max)
	}
	state.Board.XYs[a.X][a.Y].Unit = unit

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
