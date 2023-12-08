package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

var ConditionMap map[string]func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error)

func init() {
	ConditionMap = map[string]func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error){
		AlwaysFailCondition:  PassAlwaysFail,
		UnitMissingCondition: PassUnitMissing,
		ManaAboveCondition:   PassManaAbove,
		ManaBelowCondition:   PassManaBelow,
		MatchCondition:       PassMatch,
		ContainsCondition:    PassContains,
		StatAboveCondition:   PassStatAbove,
		StatBelowCondition:   PassStatBelow,
	}
}
