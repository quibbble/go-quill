package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type When string

const (
	Before When = "Before"
	After  When = "After"
)

type BuildHook func(uuid uuid.UUID, when, typ string, conditions []ICondition, event IEvent, reuse []ICondition) (IHook, error)

type IHook interface {
	UUID() uuid.UUID
	Trigger(when When, typ string) bool
	Pass(engine IEngine, state IState) (bool, error)
	Event() IEvent
	Reuse(engine IEngine, state IState) (bool, error)
}
