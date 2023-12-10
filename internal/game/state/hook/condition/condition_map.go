package condition

import (
	"context"
	"reflect"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

var ConditionMap map[string]struct {
	Type reflect.Type
	Pass func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error)
}

func init() {
	ConditionMap = map[string]struct {
		Type reflect.Type
		Pass func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error)
	}{
		ContainsCondition: {
			Type: reflect.TypeOf(&ContainsArgs{}),
			Pass: PassContains,
		},
		FailCondition: {
			Type: reflect.TypeOf(&FailArgs{}),
			Pass: PassFail,
		},
		ManaAboveCondition: {
			Type: reflect.TypeOf(&ManaAboveArgs{}),
			Pass: PassManaAbove,
		},
		ManaBelowCondition: {
			Type: reflect.TypeOf(&ManaBelowArgs{}),
			Pass: PassManaBelow,
		},
		MatchCondition: {
			Type: reflect.TypeOf(&MatchArgs{}),
			Pass: PassMatch,
		},
		StatAboveCondition: {
			Type: reflect.TypeOf(&StatAboveArgs{}),
			Pass: PassStatAbove,
		},
		StatBelowCondition: {
			Type: reflect.TypeOf(&StatBelowArgs{}),
			Pass: PassStatBelow,
		},
		UnitMissingCondition: {
			Type: reflect.TypeOf(&UnitMissingArgs{}),
			Pass: PassUnitMissing,
		},
	}
}
