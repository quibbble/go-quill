package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const OpposingPlayerChoice = "OpposingPlayer"

type OpposingPlayerArgs struct{}

func RetrieveOpposingPlayer(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	return []uuid.UUID{state.GetOpponent(state.GetTurn())}, nil
}
