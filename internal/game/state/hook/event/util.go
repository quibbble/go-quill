package event

import (
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/parse"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

func GetPlayerChoice(engine *en.Engine, state *st.State, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	return getTypeChoice(engine, state, st.PlayerUUID, raw, targets...)
}

func GetUnitChoice(engine *en.Engine, state *st.State, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	return getTypeChoice(engine, state, st.UnitUUID, raw, targets...)
}

func GetItemChoice(engine *en.Engine, state *st.State, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	return getTypeChoice(engine, state, st.ItemUUID, raw, targets...)
}

func GetSpellChoice(engine *en.Engine, state *st.State, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	return getTypeChoice(engine, state, st.SpellUUID, raw, targets...)
}

func GetTileChoice(engine *en.Engine, state *st.State, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	return getTypeChoice(engine, state, st.TileUUID, raw, targets...)
}

func GetChoice(engine *en.Engine, state *st.State, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	choose, err := ch.NewChoose(state.Gen.New(st.ChooseUUID), raw.Type, raw.Args)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	choices, err := choose.Retrieve(engine, state, targets...)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	if len(choices) != 1 {
		return uuid.Nil, errors.ErrInvalidSliceLength
	}
	return choices[0], nil
}

func getTypeChoice(engine *en.Engine, state *st.State, typ rune, raw parse.Choose, targets ...uuid.UUID) (uuid.UUID, error) {
	choice, err := GetChoice(engine, state, raw, targets...)
	if err != nil {
		return uuid.Nil, errors.Wrap(err)
	}
	if choice.Type() != typ {
		return uuid.Nil, st.ErrInvalidUUIDType(choice, typ)
	}
	return choice, nil
}
