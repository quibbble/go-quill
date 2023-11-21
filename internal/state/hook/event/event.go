package event

import (
	"fmt"

	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Event struct {
	uuid uuid.UUID

	typ    string
	args   map[string]interface{}
	affect func(engine *en.Engine, state *st.State, args map[string]interface{}) error
}

func NewEvent(typ string, args map[string]interface{}) (*Event, error) {
	affect, ok := EventMap[typ]
	if !ok {
		return nil, fmt.Errorf("'%s' is not a valid event type", typ)
	}
	return &Event{
		uuid:   uuid.New(st.EventUUID),
		typ:    typ,
		args:   args,
		affect: affect,
	}, nil
}

func (e *Event) Type() string {
	return e.typ
}

func (e *Event) Affect(engine en.IEngine, state en.IState) error {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return fmt.Errorf("failed to convert IEngine to Engine")
	}
	sta, ok := state.(*st.State)
	if !ok {
		return fmt.Errorf("failed to convert IState to State")
	}
	return e.affect(eng, sta, e.args)
}
