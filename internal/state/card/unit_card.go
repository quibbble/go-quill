package card

import (
	"math"

	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	StructureUnit = "Structure"
	CreatureUnit  = "Creature"
)

type UnitCard struct {
	*Card

	Type       string
	DamageType string
	Attack     int
	Health     int
	Cooldown   int
	Range      int
	Movement   int
	Codex      string

	// Items that apply held traits to this card
	Items []*ItemCard
}

func NewUnitCard(builders *Builders, id string, player uuid.UUID) (*UnitCard, error) {
	if len(id) == 0 || id[0] != 'U' {
		return nil, cards.ErrInvalidCardID
	}
	card, err := cards.ParseCard(id)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	unit := card.(*cards.UnitCard)
	core, err := NewCard(builders, &unit.Card, player)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &UnitCard{
		Card:       core,
		Type:       unit.Type,
		DamageType: unit.DamageType,
		Attack:     unit.Attack,
		Health:     unit.Health,
		Cooldown:   unit.Cooldown,
		Range:      unit.Range,
		Movement:   unit.Movement,
		Codex:      unit.Codex,
		Items:      make([]*ItemCard, 0),
	}, nil
}

func (u *UnitCard) AddItem(engine *en.Engine, item *ItemCard) error {
	u.Items = append(u.Items, item)
	return nil
}

func (u *UnitCard) GetAndRemoveItem(engine *en.Engine, item uuid.UUID) (*ItemCard, error) {
	idx := -1
	var itm *ItemCard
	for i, it := range u.Items {
		if it.UUID == item {
			idx = i
			itm = it
		}
	}
	if idx < 0 {
		return nil, errors.Errorf("'%s' not found on unit", item)
	}
	u.Items = append(u.Items[:idx], u.Items[idx+1:]...)
	return itm, nil
}

func (u *UnitCard) Recode(code string) error {
	tf := map[bool]string{true: "1", false: "0"}
	if len(code) != 8 {
		return errors.ErrIndexOutOfBounds
	}
	var recode string
	for i := 0; i < len(u.Codex); i++ {
		recode += tf[(u.Codex[i] == '1') != (code[i] == '1')]
	}
	u.Codex = recode
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

// CheckRange checks whether a ranged unit may attack from x1, y1 to x2, y2
func (u *UnitCard) CheckRange(x1, y1, x2, y2 int) bool {
	x := int(math.Abs(float64(x2 - x1)))
	y := int(math.Abs(float64(y2 - y1)))

	if x > u.Range || y > u.Range {
		return false
	}
	return true
}

func (c *UnitCard) Reset(build st.BuildCard) {
	card, _ := build(c.init.ID, c.Player)
	unit := card.(*UnitCard)
	unit.UUID = c.UUID
	c = unit
}
