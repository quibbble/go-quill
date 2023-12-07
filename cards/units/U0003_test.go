package units_tests

import (
	"testing"

	"github.com/quibbble/go-quill/cards/tests"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/stretchr/testify/assert"
)

func Test_U0003(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0003")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 0, 1

	// should play card
	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, game.Board.XYs[x][y].Unit.GetUUID(), uuids[0])

	u0002, _ := game.BuildCard("U0002", tests.Player2)
	game.Board.XYs[x][y+2].Unit = u0002

	// should fail cooldown check
	err = game.AttackUnit(tests.Player1, uuids[0], u0002.GetUUID())
	assert.True(t, err != nil)

	// should attack at range
	u0002Health := game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).Health
	u0003Health := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health
	game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown = 0
	if err := game.AttackUnit(tests.Player1, uuids[0], u0002.GetUUID()); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, u0003Health, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
	assert.Equal(t, u0002Health-1, game.Board.XYs[x][y+2].Unit.(*cd.UnitCard).Health)
}
