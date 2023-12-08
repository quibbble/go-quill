package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ChooseMap map[string]func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error)

func init() {
	ChooseMap = map[string]func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error){
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
