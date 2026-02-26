package game

import "slices"

type Game struct {
	Board  []string
	Player int
}

var move = map[int]string{
	1: "X",
	2: "O",
}

func NewGame() *Game {
	return &Game{
		Board:  []string{"#", "#", "#", "#", "#", "#", "#", "#", "#"},
		Player: 1,
	}
}

// MakeMove делает ход игрока по координатам x, y в диапазоне [0..2].
// Возвращает true, если ход корректен.
func (g *Game) MakeMove(x, y int) bool {
	switch {
	case x < 0 || y < 0:
		return false
	case x > 2 || y > 2:
		return false
	case g.Board[3*x+y] != "#":
		return false
	}
	g.Board[3*x+y] = move[g.Player]
	return true
}

func (g *Game) Winner() int {
	win := [][3]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},

		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},

		{0, 4, 8},
		{2, 4, 6},
	}
	players := []string{"X", "O"}
	for ind, player := range players {
		for _, variant := range win {
			var flag = true
			for _, pos := range variant {
				if g.Board[pos] != player {
					flag = false
					break
				}
			}
			if flag {
				return ind + 1
			}
		}
	}
	return 0
}

func (g *Game) NextPlayer() {
	if g.Player == 1 {
		g.Player = 2
	} else {
		g.Player = 1
	}
}

func (g *Game) Draw() bool {
	return !slices.Contains(g.Board, "#")
}
