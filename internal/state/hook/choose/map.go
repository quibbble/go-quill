package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ChooseMap map[string]func(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error)

func init() {
	ChooseMap = map[string]func(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error){
		AdjacentChoice: RetrieveAdjacent,
		BasesChoice:    RetrieveBases,
		OwnedChoice:    RetrieveOwned,
		TargetChoice:   RetrieveTarget,
		UnitsChoice:    RetrieveUnits,
		UUIDChoice:     RetrieveUUID,
	}
}
