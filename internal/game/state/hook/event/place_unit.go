package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	PlaceUnitEvent = "PlaceUnit"
)

type PlaceUnitArgs struct {
	ChoosePlayer parse.Choose
	ChooseUnit   parse.Choose
	ChooseTile   parse.Choose
}

func PlaceUnitAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a PlaceUnitArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	unitChoice, err := GetUnitChoice(engine, state, a.ChooseUnit, targets...)
	if err != nil {
		return errors.Wrap(err)
	}
	tileChoice, err := GetTileChoice(engine, state, a.ChooseTile, targets...)
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
	state.Board.XYs[tX][tY].Unit = unit

	// battle cry trait check
	battleCrys := unit.GetTraits(tr.BattleCryTrait)
	for _, battleCry := range battleCrys {
		args := battleCry.GetArgs().(tr.BattleCryArgs)
		event, err := NewEvent(state.Gen.New(st.EventUUID), args.Event.Type, args.Event.Args)
		if err != nil {
			return errors.Wrap(err)
		}
		if err := engine.Do(event, state); err != nil {
			return errors.Wrap(err)
		}
	}

	// friends/enemies trait check
	FriendsTraitCheck(engine, state)
	EnemiesTraitCheck(engine, state)

	return nil
}
