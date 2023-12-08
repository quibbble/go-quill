package choose

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const UUIDChoice = "UUID"

type UUIDArgs struct {
	UUID uuid.UUID
}

func RetrieveUUID(ctx context.Context, args interface{}, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	var c UUIDArgs
	if err := mapstructure.Decode(args, &c); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	return []uuid.UUID{c.UUID}, nil
}
