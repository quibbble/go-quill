package card

import (
	"math"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	Unit          = "Unit"
	BaseUnit      = "Base"
	StructureUnit = "Structure"
	CreatureUnit  = "Creature"
)

type UnitCard struct {
	*Card

	Type                   string
	DamageType             string
	Attack                 int
	Health                 int
	Cooldown, BaseCooldown int
	Movement, BaseMovement int
	Range                  int
	Codex                  string

	// Items that apply held traits to this card
	Items []*ItemCard
}

func (c *UnitCard) AddTrait(engine en.IEngine, trait st.ITrait) error {
	return c.Card.addTrait(engine, trait, c)
}

func (c *UnitCard) RemoveTrait(engine en.IEngine, trait uuid.UUID) error {
	return c.Card.removeTrait(engine, trait, c)
}

func (c *UnitCard) AddItem(engine *en.Engine, item *ItemCard) error {
	c.Items = append(c.Items, item)
	return nil
}

func (c *UnitCard) GetAndRemoveItem(engine *en.Engine, item uuid.UUID) (*ItemCard, error) {
	idx := -1
	var itm *ItemCard
	for i, it := range c.Items {
		if it.UUID == item {
			idx = i
			itm = it
		}
	}
	if idx < 0 {
		return nil, errors.Errorf("'%s' not found on unit", item)
	}
	c.Items = append(c.Items[:idx], c.Items[idx+1:]...)
	return itm, nil
}

func (c *UnitCard) Recode(code string) error {
	tf := map[bool]string{true: "1", false: "0"}
	if len(code) != 8 {
		return errors.ErrIndexOutOfBounds
	}
	var recode string
	for i := 0; i < len(c.Codex); i++ {
		recode += tf[(c.Codex[i] == '1') != (code[i] == '1')]
	}
	c.Codex = recode
	return nil
}

// CheckCodex checks whether the unit may move/attack from x1, y1 to x2, y2 with it's current codex
func (c *UnitCard) CheckCodex(x1, y1, x2, y2 int) bool {
	x := x2 - x1
	y := y2 - y1

	if (x < -1 || x > 1 || y < -1 || y > 1) || (x == 0 && y == 0) {
		return false
	}

	check := (x == 0 && y == 1 && c.Codex[0] == '1') || // up
		(x == 0 && y == -1 && c.Codex[1] == '1') || // down
		(x == -1 && y == 0 && c.Codex[2] == '1') || // left
		(x == 1 && y == 0 && c.Codex[3] == '1') || // right
		(x == -1 && y == 1 && c.Codex[4] == '1') || // upper-left
		(x == 1 && y == -1 && c.Codex[5] == '1') || // lower-right
		(x == -1 && y == -1 && c.Codex[6] == '1') || // lower-left
		(x == 1 && y == 1 && c.Codex[7] == '1') // upper-right

	return check
}

// CheckRange checks whether a ranged unit may attack from x1, y1 to x2, y2
func (c *UnitCard) CheckRange(x1, y1, x2, y2 int) bool {
	x := int(math.Abs(float64(x2 - x1)))
	y := int(math.Abs(float64(y2 - y1)))

	if x > c.Range || y > c.Range {
		return false
	}
	return true
}

func (c *UnitCard) Reset(build st.BuildCard) {
	card, _ := build(c.GetID(), c.Player)
	unit := card.(*UnitCard)
	unit.UUID = c.UUID
	c = unit
}
