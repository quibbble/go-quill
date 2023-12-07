package hook

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Hook struct {
	uuid       uuid.UUID
	cardUUID   uuid.UUID // the card that registered the hook
	when       en.When
	typ        string
	conditions en.Conditions
	events     []en.IEvent
	reuse      en.Conditions
}

func NewHook(uuid, cardUUID uuid.UUID, when, typ string, conditions []en.ICondition, events []en.IEvent, reuse []en.ICondition) (en.IHook, error) {
	return &Hook{
		uuid:       uuid,
		cardUUID:   cardUUID,
		when:       en.When(when),
		typ:        typ,
		conditions: conditions,
		events:     events,
		reuse:      reuse,
	}, nil
}

func (h *Hook) GetUUID() uuid.UUID {
	return h.uuid
}

func (h *Hook) GetCardUUID() uuid.UUID {
	return h.cardUUID
}

func (h *Hook) GetType() string {
	return h.typ
}

func (h *Hook) Trigger(when en.When, typ string) bool {
	return h.when == when && h.typ == typ
}

func (h *Hook) Pass(engine en.IEngine, state en.IState, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	return h.conditions.Pass(engine, state, event, targets...)
}

func (h *Hook) Events() []en.IEvent {
	return h.events
}

func (h *Hook) Reuse(engine en.IEngine, state en.IState, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	return h.reuse.Pass(engine, state, event, targets...)
}
