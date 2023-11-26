package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type BasesChoice struct {
	Player uuid.UUID
}

func (c *BasesChoice) Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error) {
	bases := make([]uuid.UUID, 0)
	for _, tile := range state.(*st.State).Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetInit().ID == "u0001" {
			bases = append(bases, tile.Unit.UUID)
		}
	}
	return bases, nil
}
