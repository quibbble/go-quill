package engine

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type IEngine interface {
	Do(event IEvent, state IState, targets ...uuid.UUID) error
	Register(hook IHook)
	DeRegister(hook IHook)
	Events() []IEvent
}

// Engine handles the core game loop logic
type Engine struct {
	// list of active hooks
	hooks []IHook

	// list of past events applied to state
	events []IEvent
}

func NewEngine(hooks ...IHook) *Engine {
	return &Engine{
		hooks:  hooks,
		events: make([]IEvent, 0),
	}
}

func (e *Engine) Do(event IEvent, state IState, targets ...uuid.UUID) error {

	if state.GameOver() {
		return nil
	}

	var err error

	for i, hook := range e.hooks {
		if hook.Trigger(Before, event.GetType()) {

			pass, err := hook.Pass(e, state, event)
			if err != nil {
				return errors.Wrap(err)
			}
			if pass {
				if err := e.Do(hook.Event(), state); err != nil {
					return errors.Wrap(err)
				}
			}

			pass, err = hook.Reuse(e, state, event)
			if err != nil {
				return errors.Wrap(err)
			}
			if !pass {
				e.hooks = append(e.hooks[:i], e.hooks[i+1:]...)
			}
		}
	}

	e.events = append(e.events, event)
	if err = event.Affect(e, state, targets...); err != nil {
		return errors.Wrap(err)
	}

	for i, hook := range e.hooks {
		if hook.Trigger(After, event.GetType()) {

			pass, err := hook.Pass(e, state, event)
			if err != nil {
				return errors.Wrap(err)
			}
			if pass {
				if err := e.Do(hook.Event(), state); err != nil {
					return errors.Wrap(err)
				}
			}

			pass, err = hook.Reuse(e, state, event)
			if err != nil {
				return errors.Wrap(err)
			}
			if !pass {
				e.hooks = append(e.hooks[:i], e.hooks[i+1:]...)
			}
		}
	}
	return nil
}

func (e *Engine) Register(hook IHook) {
	e.hooks = append(e.hooks, hook)
}

func (e *Engine) DeRegister(hook IHook) {
	for i, h := range e.hooks {
		if hook.GetUUID() == h.GetUUID() {
			e.hooks = append(e.hooks[:i], e.hooks[i+1:]...)
			return
		}
	}
}

func (e *Engine) Events() []IEvent {
	return e.events
}
