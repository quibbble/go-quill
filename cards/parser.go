package cards

import (
	"embed"
)

//go:embed items/*.yaml
//go:embed spells/*.yaml
//go:embed units/*.yaml
var fs embed.FS

func ParseCard(id string) (*Card, error)
