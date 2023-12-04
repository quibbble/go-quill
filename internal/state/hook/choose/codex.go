package choose

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const CodexChoice = "Codex"

var codexXYs = [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}, {-1, 1}, {1, -1}, {-1, -1}, {1, 1}}

type CodexArgs struct {
	UnitType string
	X, Y     int
	Codex    string
}

func RetrieveCodex(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c CodexArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	codex := make([]uuid.UUID, 0)
	for i := 0; i < len(c.Codex); i++ {
		if c.Codex[i] == '1' {
			x, y := c.X+codexXYs[i][0], c.Y+codexXYs[i][1]
			if x < 0 || x > st.Cols || y < 0 || y > st.Rows {
				continue
			}
			tile := state.(*st.State).Board.XYs[x][y]
			if tile.Unit != nil {
				unit := tile.Unit.(*cd.UnitCard)
				if (c.UnitType == cd.Unit) || (c.UnitType == unit.Type) {
					codex = append(codex, unit.UUID)
				}
			}
		}
	}
	return codex, nil
}
