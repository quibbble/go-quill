package choose

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const TilesChoice = "Tiles"

type TilesArgs struct {
	Empty bool
}

func RetrieveTiles(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c TilesArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	tiles := make([]uuid.UUID, 0)
	for _, tile := range state.(*st.State).Board.UUIDs {
		if (tile.Unit == nil) == c.Empty {
			tiles = append(tiles, tile.UUID)
		}
	}
	return tiles, nil
}