package engine

import "github.com/quibbble/go-quill/pkg/uuid"

type BuildTargetReq func(uuid uuid.UUID, typ string, args interface{}) (ITargetReq, error)

type ITargetReq interface {
	Type() string
	Validate(engine IEngine, state IState, target uuid.UUID, pior ...uuid.UUID) (bool, error)
}
