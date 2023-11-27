package cards

import (
	"embed"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/quibbble/go-quill/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidCardID = errors.Errorf("invalid card id")
)

//go:embed items/*.yaml
//go:embed spells/*.yaml
//go:embed units/*.yaml
var fs embed.FS

func ParseCard(id string) (ICard, error) {
	if len(id) == 0 {
		return nil, ErrInvalidCardID
	}

	var (
		card ICard
		raw  []byte
		err  error
	)

	switch id[0] {
	case 'I':
		card = ItemCard{}
		raw, err = fs.ReadFile(fmt.Sprintf("items/%s.yaml", id))
	case 'S':
		card = SpellCard{}
		raw, err = fs.ReadFile(fmt.Sprintf("spells/%s.yaml", id))
	case 'U':
		card = UnitCard{}
		raw, err = fs.ReadFile(fmt.Sprintf("units/%s.yaml", id))
	default:
		return nil, ErrInvalidCardID
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return nil, errors.Wrap(err)
	}
	if err := mapstructure.Decode(m, &card); err != nil {
		return nil, errors.Wrap(err)
	}
	return card, nil
}
