package trait

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type Trait struct {
	uuid uuid.UUID
	Type string
	Args interface{}

	// useful in the frontend for certain traits that are too complex to generate
	Description string

	// which item/spell/unit created the trait when not initially part of a card
	CreatedBy *uuid.UUID

	add    func(t *Trait, card st.ICard) error
	remove func(t *Trait, card st.ICard) error
}

func NewTrait(uuid uuid.UUID, createdBy *uuid.UUID, typ string, args interface{}) (st.ITrait, error) {
	ar, ok := TraitMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	decoded := reflect.New(ar.Type).Elem().Interface()
	if err := mapstructure.Decode(args, &decoded); err != nil {
		return nil, errors.Wrap(err)
	}
	return &Trait{
		uuid:      uuid,
		Type:      typ,
		Args:      decoded,
		CreatedBy: createdBy,
		add:       ar.Add,
		remove:    ar.Remove,
	}, nil
}

func (t *Trait) GetUUID() uuid.UUID {
	return t.uuid
}

func (t *Trait) GetType() string {
	return t.Type
}

func (t *Trait) GetArgs() interface{} {
	return t.Args
}

func (t *Trait) GetCreatedBy() *uuid.UUID {
	return t.CreatedBy
}

func (t *Trait) Add(card st.ICard) error {
	if t.add == nil {
		return nil
	}
	return t.add(t, card)
}

func (t *Trait) Remove(card st.ICard) error {
	if t.remove == nil {
		return nil
	}
	return t.remove(t, card)
}
