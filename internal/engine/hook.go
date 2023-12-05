package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type When string

const (
	Before When = "Before"
	After  When = "After"
)

type BuildHook func(uuid uuid.UUID, when, typ string, conditions []ICondition, event IEvent, reuse []ICondition) (IHook, error)

type IHook interface {
	GetUUID() uuid.UUID
	GetType() string
	Trigger(when When, typ string) bool
	Pass(engine IEngine, state IState, event IEvent, targets ...uuid.UUID) (bool, error)
	Event() IEvent
	Reuse(engine IEngine, state IState, event IEvent, targets ...uuid.UUID) (bool, error)
}
