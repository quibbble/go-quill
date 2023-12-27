package trait

import "github.com/quibbble/go-quill/parse"

const (
	EnrageTrait = "Enrage"
)

type EnrageArgs struct {
	Description string
	Hooks       []parse.Hook
	Events      []parse.Event
}
