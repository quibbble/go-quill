package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Condition struct {
	uuid uuid.UUID

	typ  string
	not  bool
	args interface{}
	pass func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error)
}

func NewCondition(uuid uuid.UUID, typ string, not bool, args interface{}) (en.ICondition, error) {
	pass, ok := ConditionMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	return &Condition{
		uuid: uuid,
		typ:  typ,
		not:  not,
		args: args,
		pass: pass,
	}, nil
}

func (c *Condition) GetUUID() uuid.UUID {
	return c.uuid
}

func (c *Condition) GetType() string {
	return c.typ
}

func (c *Condition) GetArgs() interface{} {
	return c.args
}

func (c *Condition) Pass(ctx context.Context, engine en.IEngine, state en.IState) (bool, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	pass, err := c.pass(ctx, c.args, eng, sta)
	if err != nil {
		return false, errors.Wrap(err)
	}
	if c.not {
		return !pass, nil
	}
	return pass, nil
}

func SliceToConditions(conditions []*Condition) en.Conditions {
	c := make([]en.ICondition, 0)
	for _, condition := range conditions {
		c = append(c, condition)
	}
	return en.Conditions(c)
}
