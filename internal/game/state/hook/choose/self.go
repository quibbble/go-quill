package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const SelfChoice = "Self"

type SelfArgs struct{}

func RetrieveSelf(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	cardCtx := ctx.Value(en.HookCardCtx)
	if cardCtx == nil {
		cardCtx := ctx.Value(en.CardCtx)
		if cardCtx == nil {
			return nil, errors.ErrMissingContext
		}
	}
	cardUUID := cardCtx.(uuid.UUID)
	return []uuid.UUID{cardUUID}, nil
}
