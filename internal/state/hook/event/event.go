package event

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Choose struct {
	Type string
	Args interface{}
}

type Trait struct {
	Type string
	Args interface{}
}

type Event struct {
	uuid uuid.UUID

	typ    string
	args   interface{}
	affect func(engine *en.Engine, state *st.State, args interface{}, targets ...uuid.UUID) error
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

func (e *Event) Type() string {
	return e.typ
}

func (e *Event) Affect(engine en.IEngine, state en.IState, targets ...uuid.UUID) error {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return e.affect(eng, sta, e.args, targets...)
}
