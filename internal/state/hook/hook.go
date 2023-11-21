package hook

import (
	en "github.com/quibbble/go-quill/internal/engine"
	cd "github.com/quibbble/go-quill/internal/state/hook/condition"
	ev "github.com/quibbble/go-quill/internal/state/hook/event"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Hook struct {
	uuid       uuid.UUID
	when       en.When
	conditions []*cd.Condition
	event      *ev.Event
	reuse      []*cd.Condition
}

func (h *Hook) UUID() uuid.UUID {
	return h.uuid
}

func (h *Hook) Trigger(when en.When, typ string) bool {
	return h.when == when && h.event.Type() == typ
}

func (h *Hook) Pass(engine en.IEngine, state en.IState) (bool, error) {
	return cd.SliceToConditions(h.conditions).Pass(engine, state)
}

func (h *Hook) Event() en.IEvent {
	return h.event
}

func (h *Hook) Reuse(engine en.IEngine, state en.IState) (bool, error) {
	return cd.SliceToConditions(h.reuse).Pass(engine, state)
}
