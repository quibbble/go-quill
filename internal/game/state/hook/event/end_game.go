package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
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
	var a EndGameArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	state.Winner = &a.Winner
	return nil
}
