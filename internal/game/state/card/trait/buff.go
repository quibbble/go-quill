package trait

import (
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/parse"

	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	BuffTrait = "Buff"
)

type BuffArgs struct {
	Stat   string
	Amount int
}

func AddBuff(t *Trait, card st.ICard) error {
	a := t.GetArgs().(*BuffArgs)
	if a.Stat == cd.CostStat {
		card.SetCost(card.GetCost() - a.Amount)
		return nil
	}

	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	switch a.Stat {
	case cd.AttackStat:
		unit.Attack += a.Amount
	case cd.HealthStat:
		unit.Health += a.Amount
	case cd.BaseCooldownStat:
		unit.BaseCooldown -= a.Amount
	case cd.BaseMovementStat:
		unit.BaseMovement += a.Amount
	case cd.RangeState:
		unit.Range += a.Amount
	default:
		return errors.Errorf("cannot buff '%s' stat", a.Stat)
	}
	return nil
}

func RemoveBuff(t *Trait, card st.ICard) error {
	a := t.GetArgs().(*BuffArgs)
	if a.Stat == cd.CostStat {
		card.SetCost(card.GetCost() + a.Amount)
		return nil
	}

	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	switch a.Stat {
	case cd.AttackStat:
		unit.Attack -= a.Amount
	case cd.HealthStat:
		initHealth := unit.GetInit().(*parse.UnitCard).Health
		if unit.Health > initHealth {
			unit.Health = initHealth
		}
	case cd.BaseCooldownStat:
		unit.BaseCooldown += a.Amount
	case cd.BaseMovementStat:
		unit.BaseMovement -= a.Amount
	case cd.RangeState:
		unit.Range -= a.Amount
	default:
		return errors.Errorf("cannot buff '%s' stat", a.Stat)
	}
	return nil
}
