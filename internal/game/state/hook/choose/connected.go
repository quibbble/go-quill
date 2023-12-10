package choose

import (
	"context"
	"slices"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const ConnectedChoice = "Connected"

type ConnectedArgs struct {
	Types          []string
	ConnectionType string
	ChooseUnit     parse.Choose
}

func RetrieveConnected(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	c := args.(*ConnectedArgs)
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), c.ChooseUnit.Type, c.ChooseUnit.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return nil, errors.ErrInvalidSliceLength
	}
	choice := choices[0]

	_, _, err = state.Board.GetUnitXY(choice)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	switch c.ConnectionType {
	case AdjacentChoice:
		choose, err = NewChoose(state.Gen.New(en.ChooseUUID), AdjacentChoice, &AdjacentArgs{
			Types:            c.Types,
			ChooseUnitOrTile: c.ChooseUnit,
		})
	case CodexChoice:
		choose, err = NewChoose(state.Gen.New(en.ChooseUUID), CodexChoice, &CodexArgs{
			Types:            c.Types,
			ChooseUnitOrTile: c.ChooseUnit,
		})
	default:
		return nil, errors.Errorf("'%s' not a valid connection type", c.ConnectionType)
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// bfs
	connected := make([]uuid.UUID, 0)
	toVist := []uuid.UUID{choice}
	for len(toVist) > 0 {
		item := toVist[0]
		toVist = toVist[1:]

		connected = append(connected, item)

		ctx := context.WithValue(context.WithValue(context.Background(), en.TargetsCtx, []uuid.UUID{item}), en.CardCtx, []uuid.UUID{item})
		conns, err := choose.Retrieve(ctx, engine, state)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		for _, conn := range conns {
			if !slices.Contains(connected, conn) {
				toVist = append(toVist, conn)
			}
		}
	}
	return connected, nil
}
