package units_tests

import (
	"testing"

	"github.com/quibbble/go-quill/cards/tests"
	"github.com/stretchr/testify/assert"
)

func Test_S0001(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0001")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u0002, _ := game.BuildCard("U0002", tests.Player2)
	game.Board.XYs[x][y].Unit = u0002

	if err := game.PlayCard(tests.Player1, uuids[0], u0002.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.True(t, game.Board.XYs[x][y].Unit == nil)
}
