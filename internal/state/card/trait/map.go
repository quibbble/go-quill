package trait

import (
	en "github.com/quibbble/go-quill/internal/engine"
	st "github.com/quibbble/go-quill/internal/state"
)

const (
	add    = "Add"
	remove = "Remove"
)

var dummy = func(engine *en.Engine, args interface{}, card st.ICard) error { return nil }

var TraitMap map[string]map[string]func(engine *en.Engine, args interface{}, card st.ICard) error

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
}
