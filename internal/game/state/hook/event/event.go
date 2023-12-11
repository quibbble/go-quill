package event

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Event struct {
	uuid uuid.UUID

	typ    string
	args   interface{}
	affect func(e *Event, ctx context.Context, engine *en.Engine, state *st.State) error
}

func NewEvent(uuid uuid.UUID, typ string, args interface{}) (en.IEvent, error) {
	a, ok := EventMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	decoded := reflect.New(a.Type).Elem().Interface()
	if err := mapstructure.Decode(args, &decoded); err != nil {
		return nil, errors.Wrap(err)
	}
	return &Event{
		uuid:   uuid,
		typ:    typ,
		args:   decoded,
		affect: a.Affect,
	}, nil
}

func (e *Event) GetUUID() uuid.UUID {
	return e.uuid
}

func (e *Event) GetType() string {
	return e.typ
}

func (e *Event) GetArgs() interface{} {
	return e.args
}

func (e *Event) Affect(ctx context.Context, engine en.IEngine, state en.IState) error {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return e.affect(e, ctx, eng, sta)
}
