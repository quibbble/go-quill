package choose

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const AdjacentChoice = "Adjacent"

var adjacentXYs = [][]int{{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}}

type AdjacentArgs struct {
	UnitType string
	X, Y     int
}

func RetrieveAdjacent(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c AdjacentArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	adjacent := make([]uuid.UUID, 0)
	for _, xy := range adjacentXYs {
		x, y := c.X+xy[0], c.Y+xy[1]
		if x < 0 || x > st.Cols || y < 0 || y > st.Rows {
			continue
		}

		tile := state.(*st.State).Board.XYs[x][y]
		if tile.Unit != nil && tile.Unit.GetID() != baseID {
			unit := tile.Unit.(*cd.UnitCard)
			if (c.UnitType == cd.Unit) || (c.UnitType == unit.Type) {
				adjacent = append(adjacent, unit.UUID)
			}
		}
	}
	return adjacent, nil
}
