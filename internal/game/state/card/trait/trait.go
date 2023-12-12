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
	typ  string
	args interface{}

	add    func(t *Trait, card st.ICard) error
	remove func(t *Trait, card st.ICard) error
}

func NewTrait(uuid uuid.UUID, typ string, args interface{}) (st.ITrait, error) {
	ar, ok := TraitMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	decoded := reflect.New(ar.Type).Elem().Interface()
	if err := mapstructure.Decode(args, &decoded); err != nil {
		return nil, errors.Wrap(err)
	}
	return &Trait{
		uuid:   uuid,
		typ:    typ,
		args:   decoded,
		add:    ar.Add,
		remove: ar.Remove,
	}, nil
}

func (t *Trait) GetUUID() uuid.UUID {
	return t.uuid
}

func (t *Trait) GetType() string {
	return t.typ
}

func (t *Trait) GetArgs() interface{} {
	return t.args
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
