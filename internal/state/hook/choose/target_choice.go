package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type TargetChoice struct {
	Index int
}

func (c *TargetChoice) Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error) {
	if c.Index < 0 || c.Index >= len(targets) {
		return nil, errors.ErrIndexOutOfBounds
	}
	return []uuid.UUID{targets[c.Index]}, nil
}
