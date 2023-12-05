package choose

import (
	en "github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

// Choices finds the intersection of all listed Choose
type Choices struct {
	Choices []IChoose
}

func NewChoices(choices ...IChoose) IChoose {
	return &Choices{
		Choices: choices,
	}
}

func (c *Choices) Retrieve(engine en.IEngine, state en.IState, targets ...uuid.UUID) ([]uuid.UUID, error) {
	lists := make([][]uuid.UUID, 0)
	for _, choice := range c.Choices {
		uuids, err := choice.Retrieve(engine, state, targets...)
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
		return uuid.Intersect(lists[0], lists[1:]...), nil
	}
}
