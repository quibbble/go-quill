package state

import (
	"github.com/quibbble/go-quill/internal/state/card"
	"github.com/quibbble/go-quill/pkg/errors"
	"github.com/quibbble/go-quill/pkg/uuid"
)

const (
	Rows = 8
	Cols = 3

	BaseID = "u0001"
)

type Tile struct {
	UUID uuid.UUID
	X, Y int
	Unit *card.UnitCard
}

func NewTile(x, y int) *Tile {
	return &Tile{
		UUID: uuid.New(TileUUID),
		X:    x,
		Y:    y,
	}
}

type Board struct {
	XYs   [Cols][Rows]*Tile
	UUIDs map[uuid.UUID]*Tile

	Sides map[uuid.UUID]int
}

func NewBoard(player1, player2 uuid.UUID) (*Board, error) {
	board := &Board{
		XYs:   [Cols][Rows]*Tile{},
		UUIDs: make(map[uuid.UUID]*Tile),
		Sides: make(map[uuid.UUID]int),
	}
	for x := 0; x < Cols; x++ {
		for y := 0; y < Rows; y++ {
			tile := NewTile(x, y)
			board.XYs[x][y] = tile
			board.UUIDs[tile.UUID] = tile
		}
	}

	for x := 0; x < Cols; x++ {
		base1, err := card.NewUnitCard(BaseID, player1)
		if err != nil {
			return nil, errors.Wrap(err)
		}
		base2, err := card.NewUnitCard(BaseID, player2)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		board.XYs[x][0].Unit = base1
		board.XYs[x][Rows-1].Unit = base2
	}

	board.Sides[player1] = 0
	board.Sides[player2] = Rows - 1

	return board, nil
}

// GetPlayableRowRange retries the range of rows the player may play a unit
func (b *Board) GetPlayableRowRange(player uuid.UUID) (int, int) {
	var min, max int
	if b.Sides[player] == 0 {
		min = 0
		max = 2

		full := true
	exit1:
		for x := 0; x < Cols; x++ {
			for y := min; y < max; y++ {
				if b.XYs[x][y].Unit == nil {
					full = false
					break exit1
				}
			}
		}
		if full {
			max++

			for x := 0; x < Cols; x++ {
				if b.XYs[x][max].Unit == nil {
					full = false
					break
				}
			}
			if full {
				max++
			}
		}
	} else {
		min = Rows - 3
		max = Rows - 1

		full := true
	exit2:
		for x := 0; x < Cols; x++ {
			for y := min; y < max; y++ {
				if b.XYs[x][y].Unit == nil {
					full = false
					break exit2
				}
			}
		}
		if full {
			min--

			for x := 0; x < Cols; x++ {
				if b.XYs[x][min].Unit == nil {
					full = false
					break
				}
			}
			if full {
				min--
			}
		}
	}
	return min, max
}
