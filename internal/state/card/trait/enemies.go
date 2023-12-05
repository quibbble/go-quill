package trait

import "github.com/quibbble/go-quill/pkg/uuid"

const (
	EnemiesTrait = "Enemies"
)

type EnemiesArgs struct {
	Choose struct {
		Type string
		Args interface{}
	}
	Trait struct {
		Type string
		Args interface{}
	}
	Current []uuid.UUID
}
