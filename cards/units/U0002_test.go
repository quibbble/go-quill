package units_tests

import (
	"testing"

	"github.com/quibbble/go-quill/cards/tests"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/stretchr/testify/assert"
)

func Test_U0002(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0002")
	if err != nil {
		t.Fatal(err)
	}

	mana := game.Mana[tests.Player1].Amount
	x, y := 0, 1

	// should play card
	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	mana -= game.Board.XYs[x][y].Unit.GetCost()
	assert.Equal(t, game.Board.XYs[x][y].Unit.GetUUID(), uuids[0])
	assert.Equal(t, mana, game.Mana[tests.Player1].Amount)

	// should fail codex check
	err = game.MoveUnit(tests.Player1, uuids[0], x+1, y)
	assert.True(t, err != nil)

	// should move unit
	y += 1
	if err := game.MoveUnit(tests.Player1, uuids[0], x, y); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Movement, 0)

	// should fail movement check
	err = game.MoveUnit(tests.Player1, uuids[0], x, y+1)
	assert.True(t, err != nil)
}
