package condition

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Condition struct {
	uuid uuid.UUID

	typ  string
	not  bool
	args interface{}
	pass func(engine *en.Engine, state *st.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error)
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

func (c *Condition) Pass(engine en.IEngine, state en.IState, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	pass, err := c.pass(eng, sta, c.args, event, targets...)
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
