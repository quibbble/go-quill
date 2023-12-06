package choose

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const SelfChoice = "Self"

type SelfArgs struct{}

func RetrieveSelf(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	if len(targets) != 1 {
		return nil, errors.ErrIndexOutOfBounds
	}
	return []uuid.UUID{targets[0]}, nil
}
