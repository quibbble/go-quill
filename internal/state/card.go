package state

import (
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ICard interface {
	GetUUID() uuid.UUID
	GetPlayer() uuid.UUID
	GetInit() cards.ICard
	Playable(engine en.IEngine, state en.IState) (bool, error)
	GetTraits(typ string) []ITrait
	AddTrait(engine en.IEngine, trait ITrait) error
	RemoveTrait(engine en.IEngine, trait uuid.UUID) error
}

type ITrait interface {
	GetUUID() uuid.UUID
	GetType() string
	GetArgs() interface{}
	Add(engine en.IEngine, card ICard) error
	Remove(engine en.IEngine, card ICard) error
}
