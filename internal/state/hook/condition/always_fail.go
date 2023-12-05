package condition

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
)

const AlwaysFailCondition = "AlwaysFail"

type AlwaysFailArgs struct{}

func PassAlwaysFail(engine *en.Engine, state *st.State, args interface{}, event ...en.IEvent) (bool, error) {
	return false, nil
}
