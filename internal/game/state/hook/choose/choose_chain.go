package choose

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	SetUnion     = "Union"
	SetIntersect = "Intersect"
)

type ChooseChain struct {
	SetFunction string
	ChooseChain []en.IChoose
}

func NewChooseChain(setFunction string, choices ...en.IChoose) en.IChoose {
	return &ChooseChain{
		SetFunction: setFunction,
		ChooseChain: choices,
	}
}

func (c *ChooseChain) Retrieve(ctx context.Context, engine en.IEngine, state en.IState) ([]uuid.UUID, error) {
	lists := make([][]uuid.UUID, 0)
	for _, choice := range c.ChooseChain {
		uuids, err := choice.Retrieve(ctx, engine, state)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		lists = append(lists, uuids)
	}
	switch len(lists) {
	case 0:
		return []uuid.UUID{}, nil
	case 1:
		return lists[0], nil
	default:
		switch c.SetFunction {
		case SetUnion:
			return uuid.Union(lists[0], lists[1:]...), nil
		case SetIntersect:
			return uuid.Intersect(lists[0], lists[1:]...), nil
		default:
			return nil, errors.Errorf("'%s' is not a supported set function", c.SetFunction)
		}
	}
}
