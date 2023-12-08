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
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	RemoveTraitFromCard = "RemoveTraitFromCard"
)

type RemoveTraitFromCardArgs struct {
	Trait      uuid.UUID
	ChooseCard parse.Choose
}

func RemoveTraitFromCardAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a RemoveTraitFromCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	choice, err := ch.GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return st.ErrNotFound(choice)
	}
	if card.RemoveTrait(engine, a.Trait); err != nil {
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
