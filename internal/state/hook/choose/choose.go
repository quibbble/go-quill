package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const baseID = "U0001"

type RawChoose struct {
	Type string
	Args interface{}
}

type IChoose interface {
	Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error)
}

type Choose struct {
	uuid uuid.UUID

	typ      string
	args     interface{}
	retrieve func(engine en.IEngine, state en.IState, args interface{}, targets ...uuid.UUID) ([]uuid.UUID, error)
}

func NewChoose(uuid uuid.UUID, typ string, args interface{}) (IChoose, error) {
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

func (c *Choose) Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error) {
	return c.retrieve(engine, state, c.args, targets...)
}
