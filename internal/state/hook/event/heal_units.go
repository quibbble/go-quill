package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	HealUnitsEvent = "HealUnits"
)

type HealUnitsArgs struct {
	Amount int
	Choose Choose
}

func HealUnitsAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a HealUnitsArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	choose, err := ch.NewChoice(a.Choose.Type, a.Choose.Args)
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
			uuid: uuid.New(st.EventUUID),
			typ:  HealUnitEvent,
			args: &HealUnitArgs{
				Amount: a.Amount,
				Choose: Choose{
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
