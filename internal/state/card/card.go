package card

import (
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ICard interface {
	GetUUID() uuid.UUID
	GetInit() cards.ICard
	Playable(engine en.IEngine, state en.IState) (bool, error)
	AddTrait(engine en.IEngine, trait ITrait) error
	RemoveTrait(engine en.IEngine, traitUUID uuid.UUID) error
}

type Card struct {
	init cards.ICard

	UUID uuid.UUID

	Owner uuid.UUID

	Cost int

	// Conditions required to play the card
	Conditions en.Conditions
	// Target requirements to play the card
	TargetReqs []en.ITargetReq

	// Hooks registered on unit play
	Hooks []en.IHook
	// Events applied on unit play
	Events []en.IEvent

	// Traits that modify this card
	Traits []ITrait
}

func NewCard(card *cards.Card, player uuid.UUID) Card

func (c *Card) GetUUID() uuid.UUID {
	return c.UUID
}

func (c *Card) GetInit() cards.ICard {
	return c.init
}

func (c *Card) Playable(engine en.IEngine, state en.IState) (bool, error) {
	return c.Conditions.Pass(engine, state)
}

func (c *Card) ValidTargets(engine en.IEngine, state en.IState, targets ...uuid.UUID) (bool, error) {
	if len(targets) != len(c.TargetReqs) {
		return false, nil
	}
	pass := true
	for i, req := range c.TargetReqs {
		p, err := req.Validate(engine, state, targets[i], targets[:i]...)
		if err != nil {
			return false, errors.Wrap(err)
		}
		pass = p && pass
	}
	return pass, nil
}

func (c *Card) AddTrait(engine en.IEngine, trait ITrait) error {
	c.Traits = append(c.Traits, trait)
	return trait.Add(engine, c)
}

func (c *Card) RemoveTrait(engine en.IEngine, traitUUID uuid.UUID) error {
	idx := -1
	var trait ITrait
	for i, t := range c.Traits {
		if t.GetUUID() == traitUUID {
			idx = i
			trait = t
		}
	}
	if idx < 0 {
		return errors.Errorf("failed to find card trait")
	}
	c.Traits = append(c.Traits[:idx], c.Traits[idx+1:]...)
	return trait.Remove(engine, c)
}
