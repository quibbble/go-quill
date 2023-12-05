package condition

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ConditionMap map[string]func(engine *en.Engine, state *state.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error)

func init() {
	ConditionMap = map[string]func(engine *en.Engine, state *state.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error){
		AlwaysFailCondition:          PassAlwaysFail,
		UnitMissingCondition:         PassUnitMissing,
		DamageUnitAppliedToCondition: PassDamageUnitAppliedTo,
	}
}
