package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	EndGameEvent = "EndGame"
)

type EndGameArgs struct {
	Winner uuid.UUID
}

func EndGameAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(EndGameArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	state.Winner = &a.Winner
	return nil
}
