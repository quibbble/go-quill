package damage

import (
	"slices"

	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	PureDamage     = "Pure"
	PhysicalDamage = "Physical"
	MagicDamage    = "Magic"
)

func Damage(unit *cd.UnitCard, amount int, typ string) (int, error) {
	if !slices.Contains([]string{PureDamage, PhysicalDamage, MagicDamage}, typ) {
		return 0, errors.Errorf("'%s' is not a valid damage type", typ)
	}
	reduction := 0
	for _, trait := range unit.Traits {
		if typ == PhysicalDamage && trait.GetType() == tr.ShieldTrait {
			reduction += trait.GetArgs().(tr.ShieldArgs).Amount
		} else if typ == MagicDamage && trait.GetType() == tr.WardTrait {
			reduction += trait.GetArgs().(tr.WardArgs).Amount
		}
	}
	damage := amount - reduction
	if damage < 0 {
		damage = 0
	}
	return damage, nil
}
