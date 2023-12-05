package choose

import (
	"slices"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const UnitsChoice = "Units"

type UnitsArgs struct {
	Players []uuid.UUID
}

func RetrieveUnits(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c UnitsArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	units := make([]uuid.UUID, 0)
	for _, tile := range state.(*st.State).Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetID() != baseID {
			unit := tile.Unit.(*cd.UnitCard)
			if slices.Contains(c.Players, unit.Player) {
				units = append(units, unit.UUID)
			}
		}
	}
	return units, nil
}
