package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const SelfChoice = "Self"

type SelfArgs struct{}

func RetrieveSelf(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	cardCtx := ctx.Value(en.HookCardCtx)
	if cardCtx == nil {
		cardCtx = ctx.Value(en.CardCtx)
		if cardCtx == nil {
			return nil, errors.ErrMissingContext
		}
	}
	cardUUID := cardCtx.(uuid.UUID)
	return []uuid.UUID{cardUUID}, nil
}
