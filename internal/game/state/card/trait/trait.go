package trait

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

type RawTrait struct {
	Type string
	Args interface{}
}

type Trait struct {
	uuid uuid.UUID
	typ  string
	args interface{}

	add    func(engine *en.Engine, args interface{}, card st.ICard) error
	remove func(engine *en.Engine, args interface{}, card st.ICard) error
}

func NewTrait(uuid uuid.UUID, typ string, args interface{}) (st.ITrait, error) {
	ar, ok := TraitMap[typ]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	a, ok := ar[add]
	if !ok {
		a = dummy
	}
	r, ok := ar[remove]
	if !ok {
		r = dummy
	}
	argsType, ok := ArgsTypeRegistry[fmt.Sprintf("trait.%sArgs", typ)]
	if !ok {
		return nil, errors.ErrMissingMapKey
	}
	typedArgs := reflect.New(argsType).Elem().Interface()
	if err := mapstructure.Decode(args, &typedArgs); err != nil {
		return nil, errors.ErrInterfaceConversion
	}
	return &Trait{
		uuid:   uuid,
		typ:    typ,
		args:   typedArgs,
		add:    a,
		remove: r,
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

func (t *Trait) SetArgs(args interface{}) {
	t.args = args
}

func (t *Trait) Add(engine en.IEngine, card st.ICard) error {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return t.add(eng, t.args, card)
}

func (t *Trait) Remove(engine en.IEngine, card st.ICard) error {
	eng, ok := engine.(*en.Engine)
	if !ok {
		return errors.ErrInterfaceConversion
	}
	return t.remove(eng, t.args, card)
}
