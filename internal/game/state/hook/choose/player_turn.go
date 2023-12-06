package choose

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const PlayerTurnChoice = "PlayerTurn"

type PlayerTurnArgs struct{}

func RetrievePlayerTurn(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{state.(*st.State).GetTurn()}, nil
}
