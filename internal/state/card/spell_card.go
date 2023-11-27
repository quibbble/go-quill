package card

import (
	"github.com/quibbble/go-quill/cards"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type SpellCard struct {
	Card
}

func NewSpellCard(id string, player uuid.UUID) (*SpellCard, error) {
	if len(id) == 0 || id[0] != 'S' {
		return nil, cards.ErrInvalidCardID
	}
	card, err := cards.ParseCard(id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	spell := card.(*cards.SpellCard)
	return &SpellCard{
		Card: NewCard(&spell.Card, player),
	}, nil
}
