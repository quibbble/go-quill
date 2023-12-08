package choose

import (
	"context"
	"slices"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const UnitsChoice = "Units"

type UnitsArgs struct {
	Types []string
}

func RetrieveUnits(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	var c UnitsArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	units := make([]uuid.UUID, 0)
	for _, tile := range state.Board.UUIDs {
		if tile.Unit != nil {
			unit := tile.Unit.(*cd.UnitCard)
			if len(c.Types) == 0 || slices.Contains(c.Types, unit.Type) {
				units = append(units, unit.UUID)
			}
		}
	}
	return units, nil
}
