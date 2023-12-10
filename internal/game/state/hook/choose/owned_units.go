package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const OwnedUnitsChoice = "OwnedUnits"

type OwnedUnitsArgs struct {
	ChoosePlayer parse.Choose
}

func RetrieveOwnedUnits(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {
	c := args.(*OwnedUnitsArgs)
	choose, err := NewChoose(state.Gen.New(en.ChooseUUID), c.ChoosePlayer.Type, c.ChoosePlayer.Args)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}
	if len(choices) != 1 || choices[0].Type() != en.PlayerUUID {
		return nil, errors.ErrInvalidSliceLength
	}
	owned := make([]uuid.UUID, 0)
	for _, tile := range state.Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetPlayer() == choices[0] {
			owned = append(owned, tile.Unit.GetUUID())
		}
	}
	return owned, nil
}
