package trait

import (
	"math/rand"
	"testing"

	st "github.com/quibbble/go-quill/internal/state"
	cd "github.com/quibbble/go-quill/internal/state/card"
	ch "github.com/quibbble/go-quill/internal/state/hook/choose"
	"github.com/quibbble/go-quill/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ModifyingTraitArgs(t *testing.T) {
	gen := uuid.NewGen(rand.New(rand.NewSource(0)))
	trait, err := NewTrait(gen.New(st.TraitUUID), FriendsTrait, &FriendsArgs{
		Choose: struct {
			Type string
			Args interface{}
		}{
			Type: ch.UUIDChoice,
			Args: &ch.UUIDArgs{
				UUID: gen.New(st.UnitUUID),
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

	args := trait.GetArgs().(*FriendsArgs)
	args.Current = []uuid.UUID{gen.New(st.UnitUUID), gen.New(st.UnitUUID), gen.New(st.UnitUUID)}

	assert.True(t, len(trait.GetArgs().(*FriendsArgs).Current) > 0)
}
