package parse

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/quibbble/go-quill/cards"
	"github.com/quibbble/go-quill/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidCardID = errors.Errorf("invalid card id")
)

var library = cards.Library

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
		card = &ItemCard{}
		raw, err = library.ReadFile(fmt.Sprintf("items/%s.yaml", id))
	case 'S':
		card = &SpellCard{}
		raw, err = library.ReadFile(fmt.Sprintf("spells/%s.yaml", id))
	case 'U':
		card = &UnitCard{}
		raw, err = library.ReadFile(fmt.Sprintf("units/%s.yaml", id))
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
	if err := mapstructure.Decode(m, card); err != nil {
		return nil, errors.Wrap(err)
	}

	// set unit card targets as they are the same no matter the unit
	if id[0] == 'U' {
		card.(*UnitCard).Card.Targets = []Choose{
			{
				Type: "Composite",
				Args: map[string]interface{}{
					"SetFunction": "Intersect",
					"Choices": []Choose{
						{
							Type: "Tiles",
							Args: map[string]interface{}{
								"Empty": true,
							},
						},
						{
							Type: "OwnedTiles",
							Args: map[string]interface{}{
								"ChoosePlayer": Choose{
									Type: "CurrentPlayer",
								},
							},
						},
					},
				},
			},
		}
	}
	return card, nil
}
