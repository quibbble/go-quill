package choose

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
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
	c, ok := ChooseMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	decoded := reflect.New(c.Type).Elem().Interface()
	if err := mapstructure.Decode(args, &decoded); err != nil {
		return nil, errors.Wrap(err)
	}
	return &Choose{
		uuid:     uuid,
		typ:      typ,
		args:     decoded,
		retrieve: c.Retrieve,
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
