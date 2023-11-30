package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type BuildEvent func(typ string, args interface{}) (IEvent, error)

type IEvent interface {
	Type() string
	Affect(engine IEngine, state IState, targets ...uuid.UUID) error
}
