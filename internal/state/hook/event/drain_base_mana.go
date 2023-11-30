package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DrainBaseManaEvent = "DrainBaseMana"
)

type DrainBaseManaArgs struct {
	Player uuid.UUID
	Amount int
}

func DrainBaseManaAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DrainBaseManaArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	state.Mana[a.Player].BaseAmount -= a.Amount
	if state.Mana[a.Player].BaseAmount < 0 {
		state.Mana[a.Player].BaseAmount = 0
	}
	return nil
}
