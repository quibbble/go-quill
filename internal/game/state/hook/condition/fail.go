package condition

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

const FailCondition = "Fail"

type FailArgs struct{}

func PassFail(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) (bool, error) {
	return false, nil
}
