package event

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	DrawCardEvent = "DrawCard"
)

type DrawCardArgs struct {
	Player uuid.UUID
}

func DrawCardAffect(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error {
	var a DrawCardArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	card, err := state.Deck[a.Player].Draw()
	if err != nil {
		return errors.Wrap(err)
	}
	state.Hand[a.Player].Add(*card)
	return nil
}
