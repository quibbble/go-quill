package hook

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Hook struct {
	uuid       uuid.UUID
	when       en.When
	typ        string
	conditions en.Conditions
	event      en.IEvent
	reuse      en.Conditions
}

func NewHook(uuid uuid.UUID, when, typ string, conditions []en.ICondition, event en.IEvent, reuse []en.ICondition) (en.IHook, error) {
	return &Hook{
		uuid:       uuid,
		when:       en.When(when),
		typ:        typ,
		conditions: conditions,
		event:      event,
		reuse:      reuse,
	}, nil
}

func (h *Hook) GetUUID() uuid.UUID {
	return h.uuid
}

func (h *Hook) GetType() string {
	return h.typ
}

func (h *Hook) Trigger(when en.When, typ string) bool {
	return h.when == when && h.event.GetType() == typ
}

func (h *Hook) Pass(engine en.IEngine, state en.IState, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	return h.conditions.Pass(engine, state, event, targets...)
}

func (h *Hook) Event() en.IEvent {
	return h.event
}

func (h *Hook) Reuse(engine en.IEngine, state en.IState, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	return h.reuse.Pass(engine, state, event, targets...)
}
