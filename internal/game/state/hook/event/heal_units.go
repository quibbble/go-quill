package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	HealUnitsEvent = "HealUnits"
)

type HealUnitsArgs struct {
	Amount      int
	ChooseUnits parse.Choose
}

func HealUnitsAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a HealUnitsArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), a.ChooseUnits.Type, a.ChooseUnits.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	for _, choice := range choices {
		if choice.Type() != st.UnitUUID {
			return st.ErrInvalidUUIDType(choice, st.UnitUUID)
		}
		event := &Event{
			uuid: state.Gen.New(st.EventUUID),
			typ:  HealUnitEvent,
			args: &HealUnitArgs{
				Amount: a.Amount,
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: &ch.UUIDArgs{
						UUID: choice,
					},
				},
			},
			affect: HealUnitAffect,
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}
	return nil
}
