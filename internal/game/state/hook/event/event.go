package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Event struct {
	uuid uuid.UUID

	typ    string
	args   interface{}
	affect func(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) error
}

func NewEvent(uuid uuid.UUID, typ string, args interface{}) (en.IEvent, error) {
	affect, ok := EventMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	return &Event{
		uuid:   uuid,
		typ:    typ,
		args:   args,
		affect: affect,
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
	return e.affect(ctx, e.args, eng, sta)
}
