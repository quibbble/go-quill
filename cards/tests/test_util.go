package tests

import (
	gm "github.com/quibbble/go-quill/internal/game"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	Seed    = int64(0)
	Player1 = uuid.UUID("P0000001")
	Player2 = uuid.UUID("P0000002")
)

func NewTestEnv(player uuid.UUID, ids ...string) (*gm.Game, []uuid.UUID, error) {

	d1 := map[string]int{
		"U0002": 30,
	}
	d2 := map[string]int{
		"U0002": 30,
	}

	game, err := gm.NewGame(Seed, Player1, Player2, d1, d2)
	if err != nil {
		return nil, nil, errors.Wrap(err)
	}

	uuids := make([]uuid.UUID, 0)

	for _, id := range ids {
		card, err := game.BuildCard(id, player, false)
		if err != nil {
			return nil, nil, errors.Wrap(err)
		}
		uuids = append(uuids, card.GetUUID())
		game.Hand[player].Add(card)
	}

	game.Mana[Player1].Amount = 8
	game.Mana[Player1].BaseAmount = 8
	game.Mana[Player2].Amount = 8
	game.Mana[Player2].BaseAmount = 8

	return game, uuids, nil
}
