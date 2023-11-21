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
	args map[string]interface{}
	pass func(engine *en.Engine, state *st.State, args map[string]interface{}) (bool, error)
}

func NewCondition(typ string, args map[string]interface{}) (*Condition, error) {
	pass, ok := ConditionMap[typ]
	if !ok {
		return nil, errors.Errorf("'%s' is not a valid condition type", typ)
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
		return false, errors.Errorf("failed to convert IEngine to Engine")
	}
	sta, ok := state.(*st.State)
	if !ok {
		return false, errors.Errorf("failed to convert IState to State")
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
