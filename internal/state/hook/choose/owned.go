package choose

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const OwnedChoice = "Owned"

type OwnedArgs struct {
	Player uuid.UUID
}

func RetrieveOwned(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error) {
	var c OwnedArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	owned := make([]uuid.UUID, 0)
	for _, tile := range state.(*st.State).Board.UUIDs {
		if tile.Unit != nil && tile.Unit.GetID() != baseID && tile.Unit.GetPlayer() == c.Player {
			owned = append(owned, tile.Unit.GetUUID())
		}
	}
	return owned, nil
}
