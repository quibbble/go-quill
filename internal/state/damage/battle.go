package damage

import (
	"github.com/quibbble/go-quill/internal/state/card"
)

func Battle(attackingUnit, defendingUnit *card.UnitCard) (int, int, error)
