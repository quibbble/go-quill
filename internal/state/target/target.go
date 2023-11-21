package target

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type TargetReq func(engine engine.IEngine, state engine.IState, target uuid.UUID, pior ...uuid.UUID) (bool, error)
