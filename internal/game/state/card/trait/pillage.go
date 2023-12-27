package trait

import "github.com/quibbble/go-quill/parse"

const (
	PillageTrait = "Pillage"
)

type PillageArgs struct {
	Description string
	Hooks       []parse.Hook
	Events      []parse.Event
}
