package event

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	CooldownEvent = "Cooldown"
)

type CooldownArgs struct {
	ChooseUnits parse.Choose
}

func CooldownAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a CooldownArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoose(state.Gen.New(en.ChooseUUID), a.ChooseUnits.Type, a.ChooseUnits.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		x, y, err := state.Board.GetUnitXY(choice)
		if err != nil {
			return errors.Wrap(err)
		}
		unit := state.Board.XYs[x][y].Unit.(*cd.UnitCard)

		// tired trait check
		if len(unit.GetTraits(tr.TiredTrait)) > 0 {
			continue
		}

		event := &Event{
			uuid: state.Gen.New(en.EventUUID),
			typ:  ModifyUnitEvent,
			args: ModifyUnitArgs{
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: unit.GetUUID(),
					},
				},
				Stat:   cd.CooldownStat,
				Amount: -1,
			},
			affect: ModifyUnitAffect,
		}
		if err := engine.Do(ctx, event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
