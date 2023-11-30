package damage

import (
	"github.com/quibbble/go-quill/cards"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
)

func Battle(attackingUnit, defendingUnit *cd.UnitCard) (int, int, error) {
	attackerDamage, err := Damage(defendingUnit, maths.MaxInt(attackingUnit.Attack, 0), attackingUnit.DamageType)
	if err != nil {
		return 0, 0, errors.Wrap(err)
	}
	// execute trait check
	if len(attackingUnit.GetTraits(tr.ExecuteTrait)) > 0 &&
		attackerDamage > 0 &&
		defendingUnit.Health < defendingUnit.GetInit().(*cards.UnitCard).Health {
		attackerDamage = defendingUnit.Health
	}
	defenderDamage, err := Damage(attackingUnit, maths.MaxInt(defendingUnit.Attack, 0), defendingUnit.DamageType)
	if err != nil {
		return 0, 0, errors.Wrap(err)
	}
	return attackerDamage, defenderDamage, nil
}
