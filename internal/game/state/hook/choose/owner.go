package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const OwnerChoice = "Owner"

type OwnerArgs struct {
	ChooseCard parse.Choose
}

func RetrieveOwner(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*OwnerArgs)
	choice, err := GetChoice(ctx, r.ChooseCard, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return nil, errors.ErrNilInterface
	}

	return []uuid.UUID{card.GetPlayer()}, nil
}
