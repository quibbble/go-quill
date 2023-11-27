package card

import (
	"github.com/quibbble/go-quill/cards"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ItemCard struct {
	Card

	// Traits applied to a unit when an item is held
	HeldTraits []ITrait
}

func NewItemCard(id string, player uuid.UUID) (*ItemCard, error) {
	if len(id) == 0 || id[0] != 'I' {
		return nil, cards.ErrInvalidCardID
	}
	card, err := cards.ParseCard(id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	item := card.(*cards.ItemCard)
	traits := make([]ITrait, 0)
	for _, trait := range item.HeldTraits {
		trait, err := NewTrait(trait.Type, trait.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		traits = append(traits, trait)
	}
	return &ItemCard{
		Card:       NewCard(&item.Card, player),
		HeldTraits: traits,
	}, nil
}
