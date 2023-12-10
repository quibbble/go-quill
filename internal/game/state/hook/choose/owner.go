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

func RetrieveOwner(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	a := args.(*OwnerArgs)
	choice, err := GetChoice(ctx, a.ChooseCard, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	card := state.GetCard(choice)
	if card == nil {
		return nil, errors.ErrNilInterface
	}

	return []uuid.UUID{card.GetPlayer()}, nil
}
