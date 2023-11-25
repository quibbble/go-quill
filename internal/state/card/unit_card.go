package card

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type UnitCard struct {
	Card

	DamageType string
	Attack     int
	Health     int
	Cooldown   int
	Range      int
	Movement   uint
	Codex      string

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

// CheckCodex checks whether the unit may move/attack from x1, y1 to x2, y2 with it's current codex
func (u *UnitCard) CheckCodex(x1, y1, x2, y2 int) bool {
	x := x2 - x1
	y := y2 - y1

	if (x < -1 || x > 1 || y < -1 || y > 1) || (x == 0 && y == 0) {
		return false
	}

	check := (x == 0 && y == 1 && u.Codex[0] == '1') || // up
		(x == 0 && y == -1 && u.Codex[1] == '1') || // down
		(x == -1 && y == 0 && u.Codex[2] == '1') || // left
		(x == 1 && y == 0 && u.Codex[3] == '1') || // right
		(x == -1 && y == 1 && u.Codex[4] == '1') || // upper-left
		(x == 1 && y == -1 && u.Codex[5] == '1') || // lower-right
		(x == -1 && y == -1 && u.Codex[6] == '1') || // lower-left
		(x == 1 && y == 1 && u.Codex[7] == '1') // upper-right

	return check
}
