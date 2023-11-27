package target

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var (
	TargeReqMap = map[string]func(engine *en.Engine, state *st.State, args interface{}, target uuid.UUID, pior ...uuid.UUID) (bool, error){
		UnitTarget:      UnitValidate,
		EmptyTileTarget: EmptyTileValidate,
	}
)
