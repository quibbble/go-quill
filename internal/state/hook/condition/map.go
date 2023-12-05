package condition

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
)

var ConditionMap map[string]func(engine *en.Engine, state *state.State, args interface{}, event ...en.IEvent) (bool, error)

func init() {
	ConditionMap = map[string]func(engine *en.Engine, state *state.State, args interface{}, event ...en.IEvent) (bool, error){
		AlwaysFailCondition:  PassAlwaysFail,
		UnitMissingCondition: PassUnitMissing,
	}
}
