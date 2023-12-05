package game

import (
	"testing"

	"github.com/quibbble/go-quill/internal/engine"
	"github.com/quibbble/go-quill/internal/state/hook"
	"github.com/quibbble/go-quill/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Compile(t *testing.T) {
	engine.NewEngine(&hook.Hook{})
}

func Test_NewGame(t *testing.T) {
	p1 := uuid.UUID("P0000001")
	p2 := uuid.UUID("P0000002")

	c1 := map[string]int{
		"U0002": 30,
	}
	c2 := map[string]int{
		"U0002": 30,
	}

	d1 := make([]string, 0)
	d2 := make([]string, 0)

	for id, count := range c1 {
		for i := 0; i < count; i++ {
			d1 = append(d1, id)
		}
	}
	for id, count := range c2 {
		for i := 0; i < count; i++ {
			d2 = append(d2, id)
		}
	}

	game, err := NewGame(0, p1, p2, d1, d2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.GetTurn(), p1)

	// hand := game.Hand[game.GetTurn()]
	// fmt.Printf("%+v", hand)
	// t.Fail()
}
