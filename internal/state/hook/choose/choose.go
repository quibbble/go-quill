package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Choose interface {
	Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error)
}

type XYChoice struct {
	Type string // tile or unit
	X, Y int
}
