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
	AddTraitToCard = "AddTraitToCard"
)

type AddTraitToCardArgs struct {
	Trait      parse.Trait
	ChooseCard parse.Choose
}

func AddTraitToCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a AddTraitToCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	trait, err := tr.NewTrait(state.Gen.New(st.ChooseUUID), a.Trait.Type, a.Trait.Args)
	if err != nil {
		return errors.Wrap(err)
	}

	choice, err := GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return st.ErrNotFound(choice)
	}
	if card.AddTrait(engine, trait); err != nil {
		return errors.Wrap(err)
	}
	if choice.Type() == st.UnitUUID && card.(*cd.UnitCard).Health <= 0 {
		// kill unit if health to low
		x, y, err := state.Board.GetUnitXY(choice)
		if err != nil {
			return errors.Wrap(err)
		}
		if state.Board.XYs[x][y].Unit.(*cd.UnitCard).Health <= 0 {
			event := &Event{
				uuid: state.Gen.New(st.EventUUID),
				typ:  KillUnitEvent,
				args: &KillUnitArgs{
					ChooseUnit: parse.Choose{
						Type: ch.UUIDChoice,
						Args: ch.UUIDArgs{
							UUID: choice,
						},
					},
				},
				affect: KillUnitAffect,
			}
			if err := engine.Do(context.Background(), event, state); err != nil {
				return errors.Wrap(err)
			}
		}
	}

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
