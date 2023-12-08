package engine

import (
	"context"

	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type BuildCondition func(uuid uuid.UUID, typ string, not bool, args interface{}) (ICondition, error)

// Conditions are requirements that must pass before some event(s) can happen
type ICondition interface {
	GetUUID() uuid.UUID
	GetType() string
	GetArgs() interface{}
	Pass(ctx context.Context, engine IEngine, state IState) (bool, error)
}

type Conditions []ICondition

func (c Conditions) Pass(ctx context.Context, engine IEngine, state IState) (bool, error) {
	pass := true
	for _, condition := range c {
		p, err := condition.Pass(ctx, engine, state)
		if err != nil {
			return false, errors.Wrap(err)
		}
		pass = p && pass
	}
	return pass, nil
}
