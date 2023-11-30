package target

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
	"github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	UnitTarget = "Unit"
)

type UnitArgs struct {
	Type string
}

func UnitValidate(engine *en.Engine, state *st.State, args interface{}, target uuid.UUID, pior ...uuid.UUID) (bool, error) {
	var v UnitArgs
	if err := mapstructure.Decode(args, &v); err != nil {
		return false, errors.ErrInterfaceConversion
	}
	x, y, err := state.Board.GetUnitXY(target)
	if err != nil {
		return false, errors.Wrap(err)
	}
	unit := state.Board.XYs[x][y].Unit
	if unit.GetUUID().Type() == st.UnitUUID && (v.Type == "" || v.Type == unit.(*card.UnitCard).Type) {
		return true, nil
	}
	return false, nil
}
