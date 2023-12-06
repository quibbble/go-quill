package damage

import (
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/maths"
)

func Battle(state *st.State, attacker, defender *cd.UnitCard) (int, int, error) {
	attackerDamage, err := Damage(defender, maths.MaxInt(attacker.Attack, 0), attacker.DamageType)
	if err != nil {
		return 0, 0, errors.Wrap(err)
	}

	// assassin trait check
	assassins := attacker.GetTraits(tr.AssassinTrait)
	if len(assassins) > 0 {
		_, aY, err := state.Board.GetUnitXY(attacker.UUID)
		if err != nil {
			return 0, 0, errors.Wrap(err)
		}
		_, dY, err := state.Board.GetUnitXY(defender.UUID)
		if err != nil {
			return 0, 0, errors.Wrap(err)
		}
		defenderSide := state.Board.Sides[defender.Player]
		if maths.AbsInt(defenderSide-aY) < maths.AbsInt(defenderSide-dY) {
			for _, trait := range assassins {
				args := trait.GetArgs().(tr.AssassinArgs)
				attackerDamage += args.Amount
			}
		}
	}

	defenderDamage, err := Damage(attacker, maths.MaxInt(defender.Attack, 0), defender.DamageType)
	if err != nil {
		return 0, 0, errors.Wrap(err)
	}

	// spiky trait check
	spikys := defender.GetTraits(tr.SpikyTrait)
	for _, spiky := range spikys {
		args := spiky.GetArgs().(tr.SpikyArgs)
		defenderDamage += args.Amount
	}

	return attackerDamage, defenderDamage, nil
}
