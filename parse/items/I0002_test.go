package units_tests

import (
	"testing"

	cd "github.com/quibbble/go-quill/internal/game/state/card"
	tr "github.com/quibbble/go-quill/internal/game/state/card/trait"
	"github.com/quibbble/go-quill/parse/tests"
	"github.com/stretchr/testify/assert"
)

func Test_I0002(t *testing.T) {
	game, uuids, err := tests.NewTestEnv(tests.Player1, "I0002")
	if err != nil {
		t.Fatal(err)
	}

	x, y := 1, 1

	u1, _ := game.BuildCard("U0002", tests.Player2, false)
	game.Board.XYs[x][y].Unit = u1

	health := u1.(*cd.UnitCard).Health

	if err := game.PlayCard(tests.Player1, uuids[0], u1.GetUUID()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, health+2, game.Board.XYs[x][y].Unit.(*cd.UnitCard).Health)
	assert.Equal(t, 1, len(u1.(*cd.UnitCard).GetTraits(tr.ShieldTrait)))
}
