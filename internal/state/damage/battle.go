package damage

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state/card"
)

func CreateBattleEvent(attackingUnit, defendingUnit *card.UnitCard) (engine.IEvent, error)
