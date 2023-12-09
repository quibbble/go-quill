package trait

import "github.com/quibbble/go-quill/parse"

const (
	DeathCryTrait = "DeathCry"
)

type DeathCryArgs struct {
	Hooks  []parse.Hook
	Events []parse.Event
}
