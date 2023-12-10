package state

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ICard interface {
	GetID() string
	GetUUID() uuid.UUID
	GetPlayer() uuid.UUID
	GetCost() int
	SetCost(cost int)
	GetInit() parse.ICard
	GetEvents() []en.IEvent
	GetHooks() []en.IHook
	Playable(engine en.IEngine, state en.IState) (bool, error)
	NextTargets(ctx context.Context, engine en.IEngine, state en.IState) ([]uuid.UUID, error)
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

type BuildTrait func(uuid uuid.UUID, typ string, args interface{}) (ITrait, error)
