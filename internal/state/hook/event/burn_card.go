package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	BurnCardEvent = "BurnCard"
)

type BurnCardArgs struct {
	Player uuid.UUID
}

func BurnCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(BurnCardArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	card, err := state.Deck[a.Player].Draw()
	if err != nil {
		return errors.Wrap(err)
	}
	state.Trash[a.Player].Add(*card)
	return nil
}
