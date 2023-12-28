package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const CardTraitsChoice = "CardTraits"

type CardTraitsArgs struct {
	ChooseCard parse.Choose
	TraitType  string
}

func RetrieveCardTraits(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	args := c.Args.(*CardTraitsArgs)
	choice, err := GetChoice(ctx, args.ChooseCard, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	card := state.GetCard(choice)
	if card == nil {
		return nil, errors.ErrNilInterface
	}
	traits := card.GetTraits(args.TraitType)
	uuids := make([]uuid.UUID, 0)
	for _, trait := range traits {
		uuids = append(uuids, trait.GetUUID())
	}
	return uuids, nil
}
