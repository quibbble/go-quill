package game

import (
	"testing"

	"github.com/quibbble/go-quill/pkg/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_NewGame(t *testing.T) {
	p1 := uuid.UUID("P0000001")
	p2 := uuid.UUID("P0000002")

	d1 := map[string]int{
		"U0002": 30,
	}
	d2 := map[string]int{
		"U0002": 30,
	}

	game, err := NewGame(0, p1, p2, d1, d2)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, game.GetTurn(), p1)
}
