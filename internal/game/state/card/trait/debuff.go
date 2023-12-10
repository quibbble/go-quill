package trait

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"

	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	DebuffTrait = "Debuff"
)

type DebuffArgs struct {
	Stat   string
	Amount int
}

func AddDebuff(engine *en.Engine, args interface{}, card st.ICard) error {
	a := args.(*DebuffArgs)
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
		unit.Health -= a.Amount
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

func RemoveDebuff(engine *en.Engine, args interface{}, card st.ICard) error {
	a := args.(*DebuffArgs)
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
