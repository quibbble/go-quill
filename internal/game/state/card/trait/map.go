package trait

import (
	"reflect"

	en "github.com/quibbble/go-quill/internal/game/engine"
	st "github.com/quibbble/go-quill/internal/game/state"
)

var TraitMap map[string]struct {
	Type   reflect.Type
	Add    func(engine *en.Engine, args interface{}, card st.ICard) error
	Remove func(engine *en.Engine, args interface{}, card st.ICard) error
}

func init() {
	TraitMap = map[string]struct {
		Type   reflect.Type
		Add    func(engine *en.Engine, args interface{}, card st.ICard) error
		Remove func(engine *en.Engine, args interface{}, card st.ICard) error
	}{
		AssassinTrait: {
			Type: reflect.TypeOf(&AssassinArgs{}),
		},
		BattleCryTrait: {
			Type: reflect.TypeOf(&BattleCryArgs{}),
		},
		BerserkTrait: {
			Type: reflect.TypeOf(&BerserkArgs{}),
		},
		BuffTrait: {
			Type:   reflect.TypeOf(&BuffArgs{}),
			Add:    AddBuff,
			Remove: RemoveBuff,
		},
		DeathCryTrait: {
			Type: reflect.TypeOf(&DeathCryArgs{}),
		},
		DebuffTrait: {
			Type:   reflect.TypeOf(&DebuffArgs{}),
			Add:    AddDebuff,
			Remove: RemoveDebuff,
		},
		EnemiesTrait: {
			Type: reflect.TypeOf(&EnemiesArgs{}),
		},
		EnrageTrait: {
			Type: reflect.TypeOf(&EnrageArgs{}),
		},
		ExecuteTrait: {
			Type: reflect.TypeOf(&ExecuteArgs{}),
		},
		FriendsTrait: {
			Type: reflect.TypeOf(&FriendsArgs{}),
		},
		GiftTrait: {
			Type: reflect.TypeOf(&GiftArgs{}),
		},
		HasteTrait: {
			Type: reflect.TypeOf(&HasteArgs{}),
		},
		LobberTrait: {
			Type: reflect.TypeOf(&LobberArgs{}),
		},
		PillageTrait: {
			Type: reflect.TypeOf(&PillageArgs{}),
		},
		PoisonTrait: {
			Type: reflect.TypeOf(&PoisonArgs{}),
		},
		PurityTrait: {
			Type: reflect.TypeOf(&PurityArgs{}),
		},
		RecodeTrait: {
			Type:   reflect.TypeOf(&RecodeArgs{}),
			Add:    AddRecode,
			Remove: RemoveRecode,
		},
		ShieldTrait: {
			Type: reflect.TypeOf(&ShieldArgs{}),
		},
		SpikyTrait: {
			Type: reflect.TypeOf(&SpikyArgs{}),
		},
		ThiefTrait: {
			Type: reflect.TypeOf(&ThiefArgs{}),
		},
		TiredTrait: {
			Type: reflect.TypeOf(&TiredArgs{}),
		},
		WardTrait: {
			Type: reflect.TypeOf(&WardArgs{}),
		},
	}
}
