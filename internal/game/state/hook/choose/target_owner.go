package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const TargetOwnerChoice = "TargetOwner"

type TargetOwnerArgs struct {
	Index int
}

func RetrieveTargetOwner(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	c := args.(*TargetOwnerArgs)
	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)
	if c.Index < 0 || c.Index >= len(targets) {
		return nil, errors.ErrIndexOutOfBounds
	}
	card := state.GetCard(targets[c.Index])
	if card == nil {
		return nil, st.ErrNotFound(targets[c.Index])
	}
	return []uuid.UUID{card.GetPlayer()}, nil
}
