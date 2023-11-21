package event

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state"
)

var (
	EventMap = map[string]func(engine *engine.Engine, state *state.State, args map[string]interface{}) error{
		"EXAMPLE": func(engine *engine.Engine, state *state.State, args map[string]interface{}) error {
			return nil
		},
	}
)
