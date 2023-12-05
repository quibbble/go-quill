package condition

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const AlwaysFailCondition = "AlwaysFail"

type AlwaysFailArgs struct{}

func PassAlwaysFail(engine *en.Engine, state *st.State, args interface{}, event en.IEvent, targets ...uuid.UUID) (bool, error) {
	return false, nil
}
