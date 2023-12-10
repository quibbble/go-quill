package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const RandomChoice = "Random"

type RandomArgs struct {
	Choose parse.Choose
}

func RetrieveRandom(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	c := args.(*RandomArgs)
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), c.Choose.Type, c.Choose.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) == 0 {
		return choices, nil
	}
	return []uuid.UUID{choices[state.Rand.Intn(len(choices))]}, nil
}
