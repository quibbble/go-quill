package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DrainBaseManaEvent = "DrainBaseMana"
)

type DrainBaseManaArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func DrainBaseManaAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DrainBaseManaArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}

	playerChoice, err := GetPlayerChoice(engine, state, a.ChoosePlayer, targets...)
	if err != nil {
		return errors.Wrap(err)
	}

	state.Mana[playerChoice].BaseAmount -= a.Amount
	if state.Mana[playerChoice].BaseAmount < 0 {
		state.Mana[playerChoice].BaseAmount = 0
	}
	return nil
}
