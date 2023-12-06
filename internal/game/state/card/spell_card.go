package card

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type SpellCard struct {
	*Card
}

func (c *SpellCard) AddTrait(engine en.IEngine, trait st.ITrait) error {
	return c.Card.addTrait(engine, trait, c)
}

func (c *SpellCard) RemoveTrait(engine en.IEngine, trait uuid.UUID) error {
	return c.Card.removeTrait(engine, trait, c)
}

func (c *SpellCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player)
	spell := card.(*SpellCard)
	spell.UUID = c.UUID
	c = spell
}
