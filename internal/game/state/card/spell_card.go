package card

import (
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type SpellCard struct {
	*Card
}

func (c *SpellCard) AddTrait(trait st.ITrait) error {
	return c.Card.addTrait(trait, c)
}

func (c *SpellCard) RemoveTrait(trait uuid.UUID) error {
	return c.Card.removeTrait(trait, c)
}

func (c *SpellCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player)
	spell := card.(*SpellCard)
	spell.UUID = c.UUID
	c = spell
}
