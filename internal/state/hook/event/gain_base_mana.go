package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	GainBaseManaEvent = "gain_base_mana"
)

type GainBaseManaArgs struct {
	Player uuid.UUID
	Amount int
}

func GainBaseManaAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(GainBaseManaArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	state.Mana[a.Player].BaseAmount += a.Amount
	return nil
}
