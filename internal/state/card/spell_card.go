package card

import (
	"github.com/quibbble/go-quill/cards"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type SpellCard struct {
	*Card
}

func NewSpellCard(builders *Builders, id string, player uuid.UUID) (*SpellCard, error) {
	if len(id) == 0 || id[0] != 'S' {
		return nil, cards.ErrInvalidCardID
	}
	card, err := cards.ParseCard(id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	spell := card.(*cards.SpellCard)
	core, err := NewCard(builders, &spell.Card, player)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &SpellCard{
		Card: core,
	}, nil
}

func (c *SpellCard) Reset(build st.BuildCard) {
	card, _ := build(c.init.ID, c.Player)
	spell := card.(*SpellCard)
	spell.UUID = c.UUID
	c = spell
}
