package choose

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const CompositeChoice = "Composite"

type CompositeArgs struct {
	SetFunction string
	Choices     []parse.Choose
}

func RetrieveComposite(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	var c CompositeArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	choices := make([]en.IChoose, 0)
	for _, ch := range c.Choices {
		choose, err := NewChoose(state.Gen.New(en.ChooseUUID), ch.Type, ch.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		choices = append(choices, choose)
	}
	ch := &Choices{
		SetFunction: c.SetFunction,
		Choices:     choices,
	}
	return ch.Retrieve(ctx, engine, state)
}
