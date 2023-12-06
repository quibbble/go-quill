package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type When string

const (
	Before When = "Before"
	After  When = "After"
)

type BuildHook func(uuid, cardUUID uuid.UUID, when, typ string, conditions []ICondition, event IEvent, reuse []ICondition) (IHook, error)

// Hooks are always registered by a card
type IHook interface {
	GetUUID() uuid.UUID
	GetCardUUID() uuid.UUID
	GetType() string
	Trigger(when When, typ string) bool
	Pass(engine IEngine, state IState, event IEvent, targets ...uuid.UUID) (bool, error)
	Event() IEvent
	Reuse(engine IEngine, state IState, event IEvent, targets ...uuid.UUID) (bool, error)
}
