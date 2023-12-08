package event

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

const SwapStatsEvent = "SwapStats"

type SwapStatsArgs struct {
	Stat        string
	ChooseCardA parse.Choose
	ChooseCardB parse.Choose
}

func SwapStatsAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a SwapStatsArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	choiceA, err := ch.GetChoice(ctx, a.ChooseCardA, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	choiceB, err := ch.GetChoice(ctx, a.ChooseCardB, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	cardA := state.GetCard(choiceA)
	if cardA == nil {
		return errors.ErrNilInterface
	}
	cardB := state.GetCard(choiceB)
	if cardB == nil {
		return errors.ErrNilInterface
	}

	if a.Stat == cd.CostStat {
		costA := cardA.GetCost()
		costB := cardB.GetCost()
		cardA.SetCost(costB)
		cardB.SetCost(costA)
		return nil
	} else if cardA.GetUUID().Type() == st.UnitUUID && cardB.GetUUID().Type() == st.UnitUUID {
		unitA := cardA.(*cd.UnitCard)
		unitB := cardB.(*cd.UnitCard)
		switch a.Stat {
		case cd.AttackStat:
			attackA := unitA.Attack
			attackB := unitB.Attack
			unitA.Attack = attackB
			unitB.Attack = attackA
		case cd.HealthStat:
			healthA := unitA.Health
			healthB := unitB.Health
			unitA.Health = healthB
			unitB.Health = healthA
		case cd.MovementStat:
			moveA := unitA.Movement
			moveB := unitB.Movement
			unitA.Movement = moveB
			unitB.Movement = moveA
		case cd.CooldownStat:
			coolA := unitA.Cooldown
			coolB := unitB.Cooldown
			unitA.Cooldown = coolB
			unitB.Cooldown = coolA
		case cd.RangeState:
			rangeA := unitA.Range
			rangeB := unitB.Range
			unitA.Range = rangeB
			unitB.Range = rangeA
		}
	}
	return errors.Errorf("'%s' cannot be swapped between '%s' and '%s'", a.Stat, choiceA, choiceB)
}
