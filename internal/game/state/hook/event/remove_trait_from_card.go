package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	RemoveTraitFromCard = "RemoveTraitFromCard"
)

type RemoveTraitFromCardArgs struct {
	ChooseTrait parse.Choose
	ChooseCard  parse.Choose
}

func RemoveTraitFromCardAffect(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error {
	a := e.GetArgs().(*RemoveTraitFromCardArgs)
	traitChoice, err := ch.GetChoice(ctx, a.ChooseTrait, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	cardChoice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card := state.GetCard(cardChoice)
	if card == nil {
		return st.ErrNotFound(cardChoice)
	}
	if card.RemoveTrait(traitChoice); err != nil {
		return errors.Wrap(err)
	}
	if cardChoice.Type() == en.UnitUUID && card.(*cd.UnitCard).Health <= 0 {
		// kill unit if health to low
		x, y, err := state.Board.GetUnitXY(cardChoice)
		if err != nil {
			return errors.Wrap(err)
		}
		if state.Board.XYs[x][y].Unit.(*cd.UnitCard).Health <= 0 {
			event, err := NewEvent(state.Gen.New(en.EventUUID), KillUnitEvent, KillUnitArgs{
				ChooseUnit: parse.Choose{
					Type: ch.UUIDChoice,
					Args: ch.UUIDArgs{
						UUID: cardChoice,
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
	FriendsTraitCheck(e, engine, state)
	EnemiesTraitCheck(e, engine, state)

	return nil
}
