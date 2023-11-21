package condition

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
)

var (
	ConditionMap = map[string]func(engine *engine.Engine, state *state.State, args map[string]interface{}) (bool, error){
		"EXAMPLE": func(engine *engine.Engine, state *state.State, args map[string]interface{}) (bool, error) {
			return false, nil
		},
	}
)
