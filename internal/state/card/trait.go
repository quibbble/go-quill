package card

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ITrait interface {
	GetUUID() uuid.UUID
	Add(engine en.IEngine, card ICard) error
	Remove(engine en.IEngine, card ICard) error
}

type Trait struct {
	uuid   uuid.UUID
	typ    string
	add    func(engine *en.Engine, card *Card) error
	remove func(engine *en.Engine, card *Card) error
}

func NewTrait(typ string, args interface{}) (ITrait, error)
