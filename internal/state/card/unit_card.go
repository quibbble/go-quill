package card

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type UnitCard struct {
	Card

	DamageType string
	Attack     uint
	Health     uint
	Cooldown   uint
	Movement   uint
	Codex      uint8

	// Items that apply held traits to this card
	Items []*ItemCard
}

func NewUnitCard(id string, player uuid.UUID) (*UnitCard, error)

func (u *UnitCard) AddItem(engine *en.Engine, item *ItemCard) error {
	for _, trait := range item.HeldTraits {
		if err := u.AddTrait(engine, trait); err != nil {
			return errors.Wrap(err)
		}
	}
	u.Items = append(u.Items, item)
	return nil
}

func (u *UnitCard) RemoveItem(engine *en.Engine, item uuid.UUID) error {
	var (
		card *ItemCard
		idx  int
	)
	for i, it := range u.Items {
		if it.UUID == item {
			card = it
			idx = i
		}
	}
	for _, trait := range card.HeldTraits {
		if err := u.RemoveTrait(engine, trait.UUID); err != nil {
			return errors.Wrap(err)
		}
	}
	u.Items = append(u.Items[:idx], u.Items[idx+1:]...)
	return nil
}
