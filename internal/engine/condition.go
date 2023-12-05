package engine

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type BuildCondition func(uuid uuid.UUID, typ string, not bool, args interface{}) (ICondition, error)

type ICondition interface {
	GetUUID() uuid.UUID
	GetType() string
	GetArgs() interface{}
	Pass(engine IEngine, state IState, event ...IEvent) (bool, error)
}

type Conditions []ICondition

func (c Conditions) Pass(engine IEngine, state IState, event ...IEvent) (bool, error) {
	pass := true
	for _, condition := range c {
		p, err := condition.Pass(engine, state, event...)
		if err != nil {
			return false, errors.Wrap(err)
		}
		pass = p && pass
	}
	return pass, nil
}
