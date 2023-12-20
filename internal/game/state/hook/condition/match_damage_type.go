package condition

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
)

const MatchDamageTypeCondition = "MatchDamageType"

type MatchDamageTypeArgs struct {
	EventContext en.Context
	DamageType   string
}

func PassMatchDamageType(c *Condition, ctx context.Context, engine *en.Engine, state *st.State) (bool, error) {
	p := c.GetArgs().(*MatchDamageTypeArgs)

	event := ctx.Value(p.EventContext).(en.IEvent)

	var a struct {
		DamageType string
	}
	if err := mapstructure.Decode(event.GetArgs(), &a); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	return a.DamageType == p.DamageType, nil
}
