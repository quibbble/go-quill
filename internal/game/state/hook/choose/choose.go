package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Choose struct {
	uuid uuid.UUID

	typ      string
	args     interface{}
	retrieve func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error)
}

func NewChoose(uuid uuid.UUID, typ string, args interface{}) (en.IChoose, error) {
	retrieve, ok := ChooseMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	return &Choose{
		uuid:     uuid,
		typ:      typ,
		args:     args,
		retrieve: retrieve,
	}, nil
}

func (c *Choose) Retrieve(ctx context.Context, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return nil, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return nil, errors.ErrInterfaceConversion
	}
	return c.retrieve(ctx, c.args, eng, sta)
}
