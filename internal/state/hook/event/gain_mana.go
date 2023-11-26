package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	GainManaEvent = "gain_mana"
)

type GainManaArgs struct {
	Player uuid.UUID
	Amount int
}

func GainManaAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(GainManaArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	state.Mana[a.Player].Amount += a.Amount
	return nil
}
