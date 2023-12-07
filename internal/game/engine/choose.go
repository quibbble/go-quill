package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type BuildChoose func(uuid uuid.UUID, typ string, args interface{}) (IChoose, error)

type IChoose interface {
	Retrieve(engine IEngine, state IState, targets ...uuid.UUID) ([]uuid.UUID, error)
}
