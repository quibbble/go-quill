package choose

import (
	"context"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const EventUnitChoice = "EventUnit"

type EventUnitArgs struct{}

func RetrieveEventUnit(ctx context.Context, args interface{}, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {

	event := ctx.Value(en.HookEventCtx).(en.IEvent)

	var a struct {
		ChooseUnit parse.Choose
	}
	if err := mapstructure.Decode(event.GetArgs(), &a); err != nil {
		return nil, errors.ErrInterfaceConversion
	}

	unitChoice, err := GetUnitChoice(ctx, a.ChooseUnit, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return []uuid.UUID{unitChoice}, nil
}
