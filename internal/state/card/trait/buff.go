package trait

import (
	"github.com/mitchellh/mapstructure"
	"github.com/quibbble/go-quill/cards"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"

	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	BuffTrait = "Buff"
)

type BuffArgs struct {
	Stat   string
	Amount int
}

func AddBuff(engine *en.Engine, args interface{}, card st.ICard) error {
	var a BuffArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	defer func() { args = &a }()

	if a.Stat == cd.CostStat {
		crd, ok := card.(*cd.Card)
		if !ok {
			return errors.ErrInterfaceConversion
		}
		crd.Cost -= a.Amount
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
	case cd.RangeState:
		unit.Range += a.Amount
	default:
		return errors.Errorf("cannot buff '%s' stat", a.Stat)
	}
	return nil
}

func RemoveBuff(engine *en.Engine, args interface{}, card st.ICard) error {
	var a BuffArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	if a.Stat == cd.CostStat {
		crd, ok := card.(*cd.Card)
		if !ok {
			return errors.ErrInterfaceConversion
		}
		crd.Cost += a.Amount
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
		initHealth := unit.GetInit().(*cards.UnitCard).Health
		if unit.Health > initHealth {
			unit.Health = initHealth
		}
	case cd.RangeState:
		unit.Range -= a.Amount
	default:
		return errors.Errorf("cannot buff '%s' stat", a.Stat)
	}
	return nil
}
