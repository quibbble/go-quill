package damage

import (
	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state/card"
)

func CreateDamageEvent(unit *card.UnitCard, amount int, typ string) (engine.IEvent, error)
