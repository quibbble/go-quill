package trait

import (
	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	cd "github.com/quibbble/go-quill/internal/game/state/card"

	"github.com/quibbble/go-quill/pkg/errors"
)

const (
	RecodeTrait = "Recode"
)

type RecodeArgs struct {
	Code string
}

func AddRecode(engine *en.Engine, args interface{}, card st.ICard) error {
	var a RecodeArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return unit.Recode(a.Code)
}

func RemoveRecode(engine *en.Engine, args interface{}, card st.ICard) error {
	var a RecodeArgs
	if err := mapstructure.Decode(args, &a); err != nil {
		return errors.ErrInterfaceConversion
	}
	unit, ok := card.(*cd.UnitCard)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return unit.Recode(a.Code)
}
