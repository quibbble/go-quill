package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	RecycleDeckEvent = "recycle_deck_event"
)

type RecycleDeckArgs struct {
	Player uuid.UUID
}

func RecycleDeckAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	a, ok := args.(RecycleDeckArgs)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	state.Deck[a.Player] = state.Discard[a.Player]
	state.Discard[a.Player] = st.NewEmptyDeck()
	state.Deck[a.Player].Shuffle()
	return nil
}
