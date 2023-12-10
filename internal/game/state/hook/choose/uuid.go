package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const UUIDChoice = "UUID"

type UUIDArgs struct {
	UUID uuid.UUID
}

func RetrieveUUID(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	c := args.(*UUIDArgs)
	return []uuid.UUID{c.UUID}, nil
}
