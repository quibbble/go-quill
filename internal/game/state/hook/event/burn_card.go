package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
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
	var a BurnCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	card, err := state.Deck[a.Player].Draw()
	if err != nil {
		return errors.Wrap(err)
	}
	state.Trash[a.Player].Add(*card)
	return nil
}
