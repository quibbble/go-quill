package trait

import "github.com/quibbble/go-quill/parse"

const (
	PillageTrait = "Pillage"
)

type PillageArgs struct {
	Hooks  []parse.Hook
	Events []parse.Event
}
