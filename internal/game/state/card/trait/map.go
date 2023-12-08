package trait

import (
	"reflect"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

const (
	add    = "Add"
	remove = "Remove"
)

var dummy = func(engine *en.Engine, args interface{}, card st.ICard) error { return nil }

var TraitMap map[string]map[string]func(engine *en.Engine, args interface{}, card st.ICard) error

var ArgsTypeRegistry = make(map[string]reflect.Type)

func init() {
	TraitMap = map[string]map[string]func(engine *en.Engine, args interface{}, card st.ICard) error{
		PoisonTrait:  {},
		BerserkTrait: {},
		RecodeTrait: {
			add:    AddRecode,
			remove: RemoveRecode,
		},
		BuffTrait: {
			add:    AddBuff,
			remove: RemoveBuff,
		},
		DebuffTrait: {
			add:    AddDebuff,
			remove: RemoveDebuff,
		},
		ExecuteTrait:   {},
		ShieldTrait:    {},
		WardTrait:      {},
		ThiefTrait:     {},
		PurityTrait:    {},
		PillageTrait:   {},
		BattleCryTrait: {},
		DeathCryTrait:  {},
		GiftTrait:      {},
		LobberTrait:    {},
		SpikyTrait:     {},
		FriendsTrait:   {},
		EnemiesTrait:   {},
		EnrageTrait:    {},
		AssassinTrait:  {},
	}

	types := []interface{}{
		PoisonArgs{},
		BerserkArgs{},
		RecodeArgs{},
		BuffArgs{},
		DebuffArgs{},
		ExecuteArgs{},
		ShieldArgs{},
		WardArgs{},
		ThiefArgs{},
		PurityArgs{},
		PillageArgs{},
		BattleCryArgs{},
		DeathCryArgs{},
		GiftArgs{},
		LobberArgs{},
		SpikyArgs{},
		FriendsArgs{},
		EnemiesArgs{},
		EnrageArgs{},
		AssassinArgs{},
	}
	for _, v := range types {
		ArgsTypeRegistry[reflect.TypeOf(v).String()] = reflect.TypeOf(v)
	}
}