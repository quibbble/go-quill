package game

import (
	"testing"

	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state/hook"
)

func Test_Compile(t *testing.T) {
	engine.NewEngine(&hook.Hook{})
}
