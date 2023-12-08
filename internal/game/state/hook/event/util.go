package event

import (
	"context"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

func GetPlayerChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, st.PlayerUUID, engine, state)
}

func GetUnitChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, st.UnitUUID, engine, state)
}

func GetItemChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, st.ItemUUID, engine, state)
}

func GetSpellChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, st.SpellUUID, engine, state)
}

func GetTileChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	return getTypeChoice(ctx, raw, st.TileUUID, engine, state)
}

func GetChoice(ctx context.Context, raw parse.Choose, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), raw.Type, raw.Args)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(ctx, engine, state)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return uuid.Nil, errors.ErrInvalidSliceLength
	}
	return choices[0], nil
}

func getTypeChoice(ctx context.Context, raw parse.Choose, typ rune, engine *en.Engine, state *st.State) (uuid.UUID, error) {
	choice, err := GetChoice(ctx, raw, engine, state)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	if choice.Type() != typ {
		return uuid.Nil, st.ErrInvalidUUIDType(choice, typ)
	}
	return choice, nil
}
