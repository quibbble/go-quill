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

const CodexChoice = "Codex"

var codexXYs = [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}, {-1, 1}, {1, -1}, {-1, -1}, {1, 1}}

type CodexArgs struct {
	Types []string
	Codex string
}

func RetrieveCodex(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c CodexArgs
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
	codex := make([]uuid.UUID, 0)
	for i := 0; i < len(c.Codex); i++ {
		if c.Codex[i] == '1' {
			x, y := x+codexXYs[i][0], y+codexXYs[i][1]
			if x < 0 || x > st.Cols || y < 0 || y > st.Rows {
				continue
			}
			tile := state.(*st.State).Board.XYs[x][y]
			if slices.Contains(c.Types, "Tile") {
				codex = append(codex, tile.UUID)
			} else if tile.Unit != nil {
				unit := tile.Unit.(*cd.UnitCard)
				if len(c.Types) == 0 || slices.Contains(c.Types, unit.Type) {
					codex = append(codex, unit.UUID)
				}
			}
		}
	}
	return codex, nil
}
