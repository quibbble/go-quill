package choose

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const TargetOwnerChoice = "TargetOwner"

type TargetOwnerArgs struct {
	Index int
}

func RetrieveTargetOwner(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c TargetArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	if c.Index < 0 || c.Index >= len(targets) {
		return nil, errors.ErrIndexOutOfBounds
	}
	card := state.(*st.State).GetCard(targets[c.Index])
	if card == nil {
		return nil, st.ErrNotFound(targets[c.Index])
	}
	return []uuid.UUID{card.GetPlayer()}, nil
}
