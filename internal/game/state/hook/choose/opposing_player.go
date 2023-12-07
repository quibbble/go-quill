package choose

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const OpposingPlayerChoice = "OpposingPlayer"

type OpposingPlayerArgs struct{}

func RetrieveOpposingPlayer(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{state.(*st.State).GetOpponent(state.(*st.State).GetTurn())}, nil
}
