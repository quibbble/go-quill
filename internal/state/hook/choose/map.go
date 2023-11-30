package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var ChooseMap = map[string]func(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error){
	BasesChoice:  RetrieveBases,
	TargetChoice: RetrieveTarget,
	UnitsChoice:  RetrieveUnits,
	UUIDChoice:   RetrieveUUID,
}
