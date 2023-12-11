package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const CompositeChoice = "Composite"

type CompositeArgs struct {
	SetFunction string
	ChooseChain []parse.Choose
}

func RetrieveComposite(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	r := c.GetArgs().(*CompositeArgs)
	choices := make([]en.IChoose, 0)
	for _, ch := range r.ChooseChain {
		choose, err := NewChoose(state.Gen.New(en.ChooseUUID), ch.Type, ch.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		choices = append(choices, choose)
	}
	ch := &ChooseChain{
		SetFunction: r.SetFunction,
		ChooseChain: choices,
	}
	return ch.Retrieve(ctx, engine, state)
}
