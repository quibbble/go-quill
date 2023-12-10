package trait

import (
	"github.com/quibbble/go-quill/parse"
)

const (
	BattleCryTrait = "BattleCry"
)

type BattleCryArgs struct {
	Hooks  []parse.Hook
	Events []parse.Event
}
