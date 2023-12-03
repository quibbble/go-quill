package engine

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type BuildCondition func(uuid uuid.UUID, typ string, args interface{}) (ICondition, error)

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
