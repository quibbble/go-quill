package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const SelfOwnerChoice = "SelfOwner"

type SelfOwnerArgs struct{}

func RetrieveSelfOwner(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {

	cardCtx := ctx.Value(en.HookCardCtx)
	if cardCtx == nil {
		cardCtx := ctx.Value(en.CardCtx)
		if cardCtx == nil {
			return nil, errors.ErrMissingContext
		}
	}

	cardUUID := cardCtx.(uuid.UUID)

	card := state.(*st.State).GetCard(cardUUID)
	if card == nil {
		return nil, st.ErrNotFound(cardUUID)
	}
	return []uuid.UUID{card.GetPlayer()}, nil
}
