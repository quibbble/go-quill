package choose

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const TargetChoice = "Target"

type TargetArgs struct {
	Index int
}

func RetrieveTarget(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	var c TargetArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)
	if c.Index < 0 || c.Index >= len(targets) {
		return nil, errors.ErrIndexOutOfBounds
	}
	return []uuid.UUID{targets[c.Index]}, nil
}
