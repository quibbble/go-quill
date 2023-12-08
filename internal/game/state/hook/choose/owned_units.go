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

const OwnedUnitsChoice = "OwnedUnits"

type OwnedUnitsArgs struct {
	ChoosePlayer parse.Choose
}

func RetrieveOwnedUnits(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	var c OwnedUnitsArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	choose, err := NewChoose(state.(*st.State).Gen.New(st.ChooseUUID), c.ChoosePlayer.Type, c.ChoosePlayer.Args)
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
	for _, tile := range state.(*st.State).Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetPlayer() == choices[0] {
			owned = append(owned, tile.Unit.GetUUID())
		}
	}
	return owned, nil
}
