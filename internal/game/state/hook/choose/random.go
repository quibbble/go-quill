package choose

import (
	"github.com/mitchellh/mapstructure"
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

func RetrieveRandom(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c RandomArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	choose, err := NewChoose(state.(*st.State).Gen.New(st.ChooseUUID), c.Choose.Type, c.Choose.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) == 0 {
		return choices, nil
	}
	return []uuid.UUID{choices[state.(*st.State).Rand.Intn(len(choices))]}, nil
}
