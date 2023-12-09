package trait

import (
	"math/rand"
	"testing"

	en "github.com/quibbble/go-quill/internal/game/engine"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	ch "github.com/quibbble/go-quill/internal/game/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ModifyingTraitArgs(t *testing.T) {
	gen := uuid.NewGen(rand.New(rand.NewSource(0)))
	trait, err := NewTrait(gen.New(en.TraitUUID), FriendsTrait, &FriendsArgs{
		ChooseUnits: struct {
			Type string
			Args interface{}
		}{
			Type: ch.UUIDChoice,
			Args: &ch.UUIDArgs{
				UUID: gen.New(en.UnitUUID),
			},
		},
		Trait: struct {
			Type string
			Args interface{}
		}{
			Type: BuffTrait,
			Args: &BuffArgs{
				Stat:   cd.AttackStat,
				Amount: 1,
			},
		},
		Current: make([]uuid.UUID, 0),
	})
	if err != nil {
		t.Fatal(err)
	}

	args := trait.GetArgs().(FriendsArgs)
	args.Current = []uuid.UUID{gen.New(en.UnitUUID), gen.New(en.UnitUUID), gen.New(en.UnitUUID)}
	trait.SetArgs(args)

	assert.True(t, len(trait.GetArgs().(FriendsArgs).Current) > 0)
}
