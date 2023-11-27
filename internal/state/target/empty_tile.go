package target

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	EmptyTileTarget = "EmptyTile"
)

type EmptyTileArgs struct{}

func EmptyTileValidate(engine *en.Engine, state *st.State, args interface{}, target uuid.UUID, pior ...uuid.UUID) (bool, error) {
	tile, ok := state.Board.UUIDs[target]
	if !ok {
		return false, errors.ErrMissingMapKey
	}
	if tile.Unit != nil {
		return true, nil
	}
	return false, nil
}
