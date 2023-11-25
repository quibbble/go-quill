package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type UUIDChoice struct {
	UUID uuid.UUID
}

func (c *UUIDChoice) Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{c.UUID}, nil
}
