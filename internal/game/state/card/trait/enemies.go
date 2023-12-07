package trait

import (
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	EnemiesTrait = "Enemies"
)

type EnemiesArgs struct {
	ChooseUnits parse.Choose
	Trait       parse.Trait
	Current     []uuid.UUID
}
