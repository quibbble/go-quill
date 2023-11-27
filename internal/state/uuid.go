package state

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	TileUUID      = 'T'
	UnitUUID      = 'U'
	SpellUUID     = 'S'
	ItemUUID      = 'I'
	PlayerUUID    = 'P'
	EngineUUID    = 'E'
	EventUUID     = 'V'
	TargetReqUUID = 'R'
	HookUUID      = 'H'
	ConditionUUID = 'C'
)

var (
	ErrInvalidUUIDType = func(uuid uuid.UUID, expectedType rune) error {
		return errors.Errorf("'%s' is not of type '%s'", string(uuid), string(expectedType))
	}
)
