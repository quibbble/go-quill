package state

import (
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

var (
	ErrNotFound = func(uuid uuid.UUID) error { return errors.Errorf("%s not found", uuid) }
)
