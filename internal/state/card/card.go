package card

import (
	"github.com/quibbble/go-quill/cards"
	"github.com/quibbble/go-quill/internal/engine"
	en "github.com/quibbble/go-quill/internal/engine"
	tr "github.com/quibbble/go-quill/internal/state/target"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Trait struct {
	UUID   uuid.UUID
	Owner  uuid.UUID
	Type   string
	Add    func(engine *engine.Engine, card *Card) error
	Remove func(engine *engine.Engine, card *Card) error
}

type Card struct {
	// data used to initialize a card
	init *cards.UnitCard

	UUID uuid.UUID

	Owner uuid.UUID

	Cost uint

	// Conditions required to play the card
	Conditions en.Conditions
	// Target requirements to play the card
	TargetReqs []tr.TargetReq

	// Events applied on unit play
	Events []en.IEvent
	// Hooks registered on unit play
	Hooks []en.IHook

	// Traits that modify this card
	Traits []Trait
}

func (c *Card) GetInit() cards.Card {
	return c.init
}

func (c *Card) Playable(engine engine.IEngine, state engine.IState) (bool, error) {
	return c.Conditions.Pass(engine, state)
}

func (c *Card) ValidTargets(engine engine.IEngine, state engine.IState, targets ...uuid.UUID) (bool, error) {
	if len(targets) != len(c.TargetReqs) {
		return false, nil
	}
	pass := true
	for i, req := range c.TargetReqs {
		p, err := req(engine, state, targets[i], targets[:i]...)
		if err != nil {
			return false, errors.Wrap(err)
		}
		pass = p && pass
	}
	return pass, nil
}

func (c *Card) AddTrait(engine *engine.Engine, trait Trait) error {
	c.Traits = append(c.Traits, trait)
	return trait.Add(engine, c)
}

func (c *Card) RemoveTrait(engine *engine.Engine, traitUUID uuid.UUID) error {
	idx := -1
	trait := Trait{}
	for i, t := range c.Traits {
		if trait.UUID == traitUUID {
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
