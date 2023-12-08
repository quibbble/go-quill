package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

const AlwaysFailCondition = "AlwaysFail"

type AlwaysFailArgs struct{}

func PassAlwaysFail(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	return false, nil
}
