package go_quill

import (
	"testing"

	bg "github.com/quibbble/go-boardgame"
	"github.com/stretchr/testify/assert"
)

func Test_Quill(t *testing.T) {
	game, err := NewQuill(&bg.BoardGameOptions{
		Teams: []string{"A", "B"},
		MoreOptions: QuillMoreOptions{
			Seed: 123,
			Decks: [][]string{
				{
					"S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001",
					"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
					"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
				},
				{
					"S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001", "S0001",
					"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
					"U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002", "U0002",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(game.teams) == 2)
}
