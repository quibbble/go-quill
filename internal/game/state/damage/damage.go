package damage

import (
	"slices"

	"github.com/mitchellh/mapstructure"
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
			var args tr.ShieldArgs
			if err := mapstructure.Decode(trait.GetArgs(), &args); err != nil {
				return 0, errors.Wrap(err)
			}
			reduction += args.Amount
		} else if typ == MagicDamage && trait.GetType() == tr.WardTrait {
			var args tr.WardArgs
			if err := mapstructure.Decode(trait.GetArgs(), &args); err != nil {
				return 0, errors.Wrap(err)
			}
			reduction += args.Amount
		}
	}
	damage := amount - reduction
	if damage < 0 {
		damage = 0
	}
	return damage, nil
}
