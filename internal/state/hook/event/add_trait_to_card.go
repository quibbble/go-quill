package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	tr "github.com/quibbble/go-quill/internal/state/card/trait"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	AddTraitToCard = "AddTraitToCard"
)

type AddTraitToCardArgs struct {
	Trait  Trait
	Choose Choose
}

func AddTraitToCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a AddTraitToCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	trait, err := tr.NewTrait(state.Gen.New(st.ChooseUUID), a.Trait.Type, a.Trait.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), a.Choose.Type, a.Choose.Args)
	if err != nil {
		return errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	if len(choices) != 1 {
		return errors.ErrInvalidSliceLength
	}
	card := state.GetCard(choices[0])
	if card == nil {
		return st.ErrNotFound(choices[0])
	}
	if card.AddTrait(engine, trait); err != nil {
		return errors.Wrap(err)
	}
	if choices[0].Type() == st.UnitUUID && card.(*cd.UnitCard).Health <= 0 {
		// kill unit if health to low
		x, y, err := state.Board.GetUnitXY(choices[0])
		if err != nil {
			return errors.Wrap(err)
		}
		if state.Board.XYs[x][y].Unit.(*cd.UnitCard).Health <= 0 {
			event := &Event{
				uuid: state.Gen.New(st.EventUUID),
				typ:  KillUnitEvent,
				args: &KillUnitArgs{
					Choose: Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: choices[0],
						},
					},
				},
				affect: KillUnitAffect,
			}
			if err := engine.Do(event, state); err != nil {
				return errors.Wrap(err)
			}
		}
	}

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
