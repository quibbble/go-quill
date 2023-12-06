package target

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type TargetReq struct {
	uuid uuid.UUID

	typ      string
	args     interface{}
	validate func(engine *en.Engine, state *st.State, args interface{}, target uuid.UUID, pior ...uuid.UUID) (bool, error)
}

func NewTargetReq(uuid uuid.UUID, typ string, args interface{}) (en.ITargetReq, error) {
	validate, ok := TargeReqMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	return &TargetReq{
		uuid:     uuid,
		typ:      typ,
		args:     args,
		validate: validate,
	}, nil
}

func (t *TargetReq) Type() string {
	return t.typ
}

func (t *TargetReq) Validate(engine en.IEngine, state en.IState, target uuid.UUID, pior ...uuid.UUID) (bool, error) {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	sta, ok := state.(*st.State)
	if !ok {
		return false, errors.ErrInterfaceConversion
	}
	return t.validate(eng, sta, t.args, target, pior...)
}
