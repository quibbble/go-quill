package units_tests

import (
	"testing"

	"github.com/quibbble/go-quill/cards/tests"
	cd "github.com/quibbble/go-quill/internal/game/state/card"
	"github.com/quibbble/go-quill/pkg/maths"
	"github.com/stretchr/testify/assert"
)

func Test_S0008(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "S0008")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2)
	game.Board.XYs[x][y].Unit = u1

	cooldown := game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, maths.MaxInt(0, cooldown-2), game.Board.XYs[x][y].Unit.(*cd.UnitCard).Cooldown)
}
