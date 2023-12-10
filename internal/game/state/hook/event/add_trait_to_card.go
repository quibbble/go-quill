package event

import (
	"context"

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
	a := args.(*AddTraitToCardArgs)
	trait, err := tr.NewTrait(state.Gen.New(en.ChooseUUID), a.Trait.Type, a.Trait.Args)
	if err != nil {
		return errors.Wrap(err)
	}

	choice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
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
	if choice.Type() == en.UnitUUID && card.(*cd.UnitCard).Health <= 0 {
		// kill unit if health to low
		x, y, err := state.Board.GetUnitXY(choice)
		if err != nil {
			return errors.Wrap(err)
		}
		if state.Board.XYs[x][y].Unit.(*cd.UnitCard).Health <= 0 {
			event, err := NewEvent(state.Gen.New(en.EventUUID), KillUnitEvent, KillUnitArgs{
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: choice,
					},
				},
			})
			if err != nil {
				return errors.Wrap(err)
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
