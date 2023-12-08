package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const CurrentPlayerChoice = "CurrentPlayer"

type CurrentPlayerArgs struct{}

func RetrieveCurrentPlayer(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	return []uuid.UUID{state.GetTurn()}, nil
}
