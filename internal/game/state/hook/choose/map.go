package choose

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ChooseMap map[string]func(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error)

func init() {
	ChooseMap = map[string]func(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error){
		AdjacentChoice:       RetrieveAdjacent,
		CodexChoice:          RetrieveCodex,
		OwnedUnitsChoice:     RetrieveOwnedUnits,
		OwnedTilesChoice:     RetrieveOwnedTiles,
		TargetChoice:         RetrieveTarget,
		UnitsChoice:          RetrieveUnits,
		TilesChoice:          RetrieveTiles,
		UUIDChoice:           RetrieveUUID,
		SelfChoice:           RetrieveSelf,
		ConnectedChoice:      RetrieveConnected,
		CompositeChoice:      RetrieveComposite,
		CurrentPlayerChoice:  RetrieveCurrentPlayer,
		OpposingPlayerChoice: RetrieveOpposingPlayer,
	}
}
