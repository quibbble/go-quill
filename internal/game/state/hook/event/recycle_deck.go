package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	RecycleDeckEvent = "RecycleDeck"
)

type RecycleDeckArgs struct {
	Player uuid.UUID
}

func RecycleDeckAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a RecycleDeckArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	state.Deck[a.Player] = state.Discard[a.Player]
	state.Discard[a.Player] = st.NewEmptyDeck(state.Seed)
	state.Deck[a.Player].Shuffle()
	return nil
}
