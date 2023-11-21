package engine

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type When string

const (
	Before When = "Before"
	After  When = "After"
)

type IEngine interface {
	Do(event IEvent, state IState) error
	Register(hook IHook)
	DeRegister(hook IHook)
}

type IState interface {
}

type IHook interface {
	UUID() uuid.UUID
	Trigger(when When, typ string) bool
	Pass(engine IEngine, state IState) (bool, error)
	Event() IEvent
	Reuse(engine IEngine, state IState) (bool, error)
}

type IEvent interface {
	Type() string
	Affect(engine IEngine, state IState) error
}

type ICondition interface {
	Type() string
	Pass(engine IEngine, state IState) (bool, error)
}

type Conditions []ICondition

func (c Conditions) Pass(engine IEngine, state IState) (bool, error) {
	pass := true
	for _, condition := range c {
		p, err := condition.Pass(engine, state)
		if err != nil {
			return false, errors.Wrap(err)
		}
		pass = p && pass
	}
	return pass, nil
}
