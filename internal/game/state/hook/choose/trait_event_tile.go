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

const TraitEventTileChoice = "TraitEventTile"

type TraitEventTileArgs struct{}

func RetrieveTraitEventTile(c *Choose, ctx context.Context, engine *en.Engine, state *st.State) ([]uuid.UUID, error) {

	event := ctx.Value(en.TraitEventCtx).(en.IEvent)

	var a struct {
		ChooseTile parse.Choose
	}
	if err := mapstructure.Decode(event.GetArgs(), &a); err != nil {
		return nil, errors.ErrInterfaceConversion
	}

	tileChoice, err := GetTileChoice(ctx, a.ChooseTile, engine, state)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return []uuid.UUID{tileChoice}, nil
}
