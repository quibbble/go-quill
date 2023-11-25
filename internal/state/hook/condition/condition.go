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
	args interface{}
	pass func(engine *en.Engine, state *st.State, args interface{}) (bool, error)
}

func NewCondition(typ string, args interface{}) (*Condition, error) {
	pass, ok := ConditionMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	return &Condition{
		uuid: uuid.New(st.ConditionUUID),
		typ:  typ,
		args: args,
		pass: pass,
	}, nil
}

func (c *Condition) Type() string {
	return c.typ
}

func (c *Condition) Pass(engine en.IEngine, state en.IState) (bool, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	return c.pass(eng, sta, c.args)
}

func SliceToConditions(conditions []*Condition) en.Conditions {
	c := make([]en.ICondition, 0)
	for _, condition := range conditions {
		c = append(c, condition)
	}
	return en.Conditions(c)
}
