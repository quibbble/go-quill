package choose

import (
	"slices"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const AdjacentChoice = "Adjacent"

var adjacentXYs = [][]int{{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {1, 1}, {0, 1}, {-1, 1}, {-1, 0}}

type AdjacentArgs struct {
	Types []string
}

func RetrieveAdjacent(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c AdjacentArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	if len(targets) != 1 {
		return nil, errors.ErrInvalidSliceLength
	}
	x, y, err := state.(*st.State).Board.GetUnitXY(targets[0])
	if err != nil {
		return nil, errors.Wrap(err)
	}
	adjacent := make([]uuid.UUID, 0)
	for _, xy := range adjacentXYs {
		x, y := x+xy[0], y+xy[1]
		if x < 0 || x >= st.Cols || y < 0 || y >= st.Rows {
			continue
		}

		tile := state.(*st.State).Board.XYs[x][y]
		if slices.Contains(c.Types, "Tile") {
			adjacent = append(adjacent, tile.UUID)
		} else if tile.Unit != nil {
			unit := tile.Unit.(*cd.UnitCard)
			if len(c.Types) == 0 || slices.Contains(c.Types, unit.Type) {
				adjacent = append(adjacent, unit.UUID)
			}
		}
	}
	return adjacent, nil
}
