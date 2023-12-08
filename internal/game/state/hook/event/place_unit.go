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
	PlaceUnitEvent = "PlaceUnit"
)

type PlaceUnitArgs struct {
	ChoosePlayer parse.Choose
	ChooseUnit   parse.Choose
	ChooseTile   parse.Choose
}

func PlaceUnitAffect(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error {
	var a PlaceUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := ch.GetPlayerChoice(ctx, a.ChoosePlayer, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := ch.GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := ch.GetTileChoice(ctx, a.ChooseTile, engine, state)
	if err != nil {
		return errors.Wrap(err)
	}

	card, err := state.Hand[playerChoice].GetCard(unitChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	unit := card.(*cd.UnitCard)

	tX, tY, err := state.Board.GetTileXY(tileChoice)
	if err != nil {
		return errors.Wrap(err)
	}
	if state.Board.XYs[tX][tY].Unit != nil {
		return errors.Errorf("unit '%s' cannot be placed on a full tile", unit.UUID)
	}
	min, max := state.Board.GetPlayableRowRange(playerChoice)
	if tY < min || tY > max {
		return errors.Errorf("unit '%s' must be placed within rows %d to %d", unit.UUID, min, max)
	}
	if err := state.Hand[playerChoice].RemoveCard(unitChoice); err != nil {
		return errors.Wrap(err)
	}

	// haste trait check
	if len(unit.GetTraits(tr.HasteTrait)) > 0 {
		unit.Cooldown = 0
	}

	state.Board.XYs[tX][tY].Unit = unit

	// battle cry trait check
	battleCrys := unit.GetTraits(tr.BattleCryTrait)
	for _, battleCry := range battleCrys {
		args := battleCry.GetArgs().(tr.BattleCryArgs)
		event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(context.Background(), event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
