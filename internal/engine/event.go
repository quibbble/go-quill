package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type BuildEvent func(uuid uuid.UUID, typ string, args interface{}) (IEvent, error)

type IEvent interface {
	GetUUID() uuid.UUID
	GetType() string
	GetArgs() interface{}
	Affect(engine IEngine, state IState, targets ...uuid.UUID) error
}
