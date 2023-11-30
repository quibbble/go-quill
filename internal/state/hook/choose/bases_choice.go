package choose

import (
	"slices"

	"github.com/mitchellh/mapstructure"
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const BasesChoice = "Bases"

type BasesArgs struct {
	Players []uuid.UUID
}

func RetrieveBases(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c BasesArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	bases := make([]uuid.UUID, 0)
	for _, tile := range state.(*st.State).Board.UUIDs {
		if tile.Unit != nil {
			unit := tile.Unit.(*cd.UnitCard)
			if slices.Contains(c.Players, unit.Player) && unit.GetInit().(*cards.UnitCard).ID == "U0001" {
				bases = append(bases, unit.UUID)
			}
		}

	}
	return bases, nil
}
