package units_tests

import (
	"testing"

	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_U0008(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "U0008", "S0001", "S0010")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 2

	u0002, _ := game.BuildCard("U0002", tests.Player1, false)
	game.Board.XYs[x][y-1].Unit = u0002

	if err := game.PlayCard(tests.Player1, uuids[0], game.Board.XYs[x][y].UUID); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(game.Hooks()))

	attack := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack
	health := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health

	if err := game.PlayCard(tests.Player1, uuids[1], u0002.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, attack+1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Attack)
	assert.Equal(t, health+1, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)

	game.Mana[tests.Player1].Amount += 10

	if err := game.PlayCard(tests.Player1, uuids[2], uuids[0]); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(game.Hooks()))
}
