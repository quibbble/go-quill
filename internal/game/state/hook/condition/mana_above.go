package condition

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ev "github.com/quibbble/go-quill/internal/game/state/hook/event"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const ManaAboveCondition = "ManaAbove"

type ManaAboveArgs struct {
	ChoosePlayer parse.Choose
	Amount       int
}

func PassManaAbove(engine *en.Engine, state *st.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	var p ManaAboveArgs
	if err := mapstructure.Decode(args, &p); err != nil {
		return false, errors.ErrInterfaceConversion
	}

	playerChoice, err := ev.GetPlayerChoice(engine, state, p.ChoosePlayer, targets...)
	if err != nil {
		return false, errors.Wrap(err)
	}

	return state.Mana[playerChoice].Amount > p.Amount, nil
}
