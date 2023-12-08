package choose

import (
	"context"
	"slices"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const ConnectedChoice = "Connected"

type ConnectedArgs struct {
	Types          []string
	ConnectionType string
}

func RetrieveConnected(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	var c ConnectedArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	targets := ctx.Value(en.TargetsCtx).([]uuid.UUID)
	if len(targets) != 1 {
		return nil, errors.ErrInvalidSliceLength
	}
	_, _, err := state.(*st.State).Board.GetUnitXY(targets[0])
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var choose en.IChoose
	switch c.ConnectionType {
	case AdjacentChoice:
		choose, err = NewChoose(state.(*st.State).Gen.New(st.ChooseUUID), AdjacentChoice, &AdjacentArgs{
			Types: c.Types,
		})
	case CodexChoice:
		choose, err = NewChoose(state.(*st.State).Gen.New(st.ChooseUUID), CodexChoice, &CodexArgs{
			Types: c.Types,
		})
	default:
		return nil, errors.Errorf("'%s' not a valid connection type", c.ConnectionType)
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}

	// bfs
	connected := make([]uuid.UUID, 0)
	toVist := []uuid.UUID{targets[0]}
	for len(toVist) > 0 {
		item := toVist[0]
		toVist = toVist[1:]

		connected = append(connected, item)

		ctx := context.WithValue(context.Background(), en.TargetsCtx, []uuid.UUID{item})
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
