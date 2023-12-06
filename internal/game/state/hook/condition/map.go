package condition

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/internal/game/state"
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
