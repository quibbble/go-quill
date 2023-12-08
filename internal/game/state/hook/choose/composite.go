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
	Choices []parse.Choose
}

func RetrieveComposite(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	var c CompositeArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	choices := make([]en.IChoose, 0)
	for _, ch := range c.Choices {
		choose, err := NewChoose(state.(*st.State).Gen.New(st.ChooseUUID), ch.Type, ch.Args)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		choices = append(choices, choose)
	}
	ch := &Choices{
		Choices: choices,
	}
	return ch.Retrieve(ctx, engine, state)
}
