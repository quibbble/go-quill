package choose

import (
	"slices"

	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type UnitsChoice struct {
	Players []uuid.UUID
}

func (c *UnitsChoice) Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error) {
	units := make([]uuid.UUID, 0)
	for _, tile := range state.(*st.State).Board.UUIDs {
		if tile.Unit != nil &&
			slices.Contains(c.Players, tile.Unit.Owner) {
			units = append(units, tile.Unit.UUID)
		}
	}
	return units, nil
}
