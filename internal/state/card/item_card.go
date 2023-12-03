package card

import (
	"github.com/quibbble/go-quill/cards"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type ItemCard struct {
	*Card

	// Traits applied to a unit when an item is held
	HeldTraits []st.ITrait
}

func NewItemCard(builders *Builders, id string, player uuid.UUID) (*ItemCard, error) {
	if len(id) == 0 || id[0] != 'I' {
		return nil, cards.ErrInvalidCardID
	}
	card, err := cards.ParseCard(id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	item := card.(*cards.ItemCard)
	traits := make([]st.ITrait, 0)
	for _, trait := range item.HeldTraits {
		trait, err := builders.BuildTrait(builders.Gen.New(st.TraitUUID), trait.Type, trait.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		traits = append(traits, trait)
	}
	core, err := NewCard(builders, &item.Card, player)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &ItemCard{
		Card:       core,
		HeldTraits: traits,
	}, nil
}

func (c *ItemCard) Reset(build st.BuildCard) {
	card, _ := build(c.init.ID, c.Player)
	item := card.(*ItemCard)
	item.UUID = c.UUID
	c = item
}
