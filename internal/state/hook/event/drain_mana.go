package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DrainManaEvent = "drain_mana"
)

type DrainManaArgs struct {
	Player uuid.UUID
	Amount int
}

func DrainManaAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(DrainManaArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	state.Mana[a.Player].Amount -= a.Amount
	if state.Mana[a.Player].Amount < 0 {
		state.Mana[a.Player].Amount = 0
	}
	return nil
}
