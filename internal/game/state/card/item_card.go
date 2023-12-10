package card

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ItemCard struct {
	*Card

	// UnitCard that is holding this item
	Holder *uuid.UUID

	// Traits applied to a unit when an item is held
	HeldTraits []st.ITrait
}

func (c *ItemCard) AddTrait(engine en.IEngine, trait st.ITrait) error {
	return c.Card.addTrait(engine, trait, c)
}

func (c *ItemCard) RemoveTrait(engine en.IEngine, trait uuid.UUID) error {
	return c.Card.removeTrait(engine, trait, c)
}

func (c *ItemCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player)
	item := card.(*ItemCard)
	item.UUID = c.UUID
	c = item
}
