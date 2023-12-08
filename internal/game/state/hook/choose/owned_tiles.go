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

const OwnedTilesChoice = "OwnedTiles"

type OwnedTilesArgs struct {
	ChoosePlayer parse.Choose
}

func RetrieveOwnedTiles(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	var c OwnedTilesArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	choose, err := NewChoose(state.Gen.New(st.ChooseUUID), c.ChoosePlayer.Type, c.ChoosePlayer.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) != 1 || choices[0].Type() != st.PlayerUUID {
		return nil, errors.ErrInvalidSliceLength
	}
	owned := make([]uuid.UUID, 0)
	min, max := state.Board.GetPlayableRowRange(choices[0])
	for _, col := range state.Board.XYs {
		for y, tile := range col {
			if min <= y && y <= max {
				owned = append(owned, tile.UUID)
			}
		}
	}
	return owned, nil
}
