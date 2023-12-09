package condition

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const StatAboveCondition = "StatAbove"

type StatAboveArgs struct {
	Stat       string
	Amount     int
	ChooseCard parse.Choose
}

func PassStatAbove(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	var p StatAboveArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	choice, err := ch.GetChoice(ctx, p.ChooseCard, engine, state)
	if err != nil {
		return false, errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return false, errors.ErrNilInterface
	}

	if p.Stat == cd.CostStat {
		return p.Amount < card.GetCost(), nil
	} else if choice.Type() == en.UnitUUID {
		unit := card.(*cd.UnitCard)
		switch p.Stat {
		case cd.AttackStat:
			return p.Amount < unit.Attack, nil
		case cd.HealthStat:
			return p.Amount < unit.Health, nil
		case cd.CooldownStat:
			return p.Amount < unit.Cooldown, nil
		case cd.MovementStat:
			return p.Amount < unit.Movement, nil
		case cd.RangeState:
			return p.Amount < unit.Range, nil
		}
	}
	return false, errors.Errorf("'%s' is not a valid stat", p.Stat)
}
