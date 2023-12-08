package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ChooseMap map[string]func(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error)

func init() {
	ChooseMap = map[string]func(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error){
		AdjacentChoice:       RetrieveAdjacent,
		CodexChoice:          RetrieveCodex,
		OwnedUnitsChoice:     RetrieveOwnedUnits,
		OwnedTilesChoice:     RetrieveOwnedTiles,
		TargetChoice:         RetrieveTarget,
		UnitsChoice:          RetrieveUnits,
		TilesChoice:          RetrieveTiles,
		UUIDChoice:           RetrieveUUID,
		SelfChoice:           RetrieveSelf,
		SelfOwnerChoice:      RetrieveSelfOwner,
		ConnectedChoice:      RetrieveConnected,
		CompositeChoice:      RetrieveComposite,
		CurrentPlayerChoice:  RetrieveCurrentPlayer,
		OpposingPlayerChoice: RetrieveOpposingPlayer,
		TargetOwnerChoice:    RetrieveTargetOwner,
		RandomChoice:         RetrieveRandom,
	}
}
