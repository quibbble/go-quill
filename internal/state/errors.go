package state

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var (
	ErrUnitNotFound = func(unit uuid.UUID) error { return errors.Errorf("unit %s not found", unit) }
)
