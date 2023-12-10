package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const TilesChoice = "Tiles"

type TilesArgs struct {
	Empty bool
}

func RetrieveTiles(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	c := args.(*TilesArgs)
	tiles := make([]uuid.UUID, 0)
	for _, tile := range state.Board.UUIDs {
		if (tile.Unit == nil) == c.Empty {
			tiles = append(tiles, tile.UUID)
		}
	}
	return tiles, nil
}
