package game_test

import (
	"testing"

	"github.com/Dzetner/tic-tac-toe-grpc/game"
)

func TestGame_MakeMove(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		prepareMove bool
		x           int
		y           int
		want        bool
	}{
		{
			name:        "valid",
			prepareMove: false,
			x:           1,
			y:           1,
			want:        true,
		},
		{
			name:        "invalid x",
			prepareMove: false,
			x:           0,
			y:           1,
			want:        false,
		},
		{
			name:        "invalid y",
			prepareMove: false,
			x:           1,
			y:           0,
			want:        false,
		},
		{
			name:        "already occupied",
			prepareMove: true,
			x:           1,
			y:           1,
			want:        false,
		},

		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGame()
			if tt.prepareMove {
				ok := g.MakeMove(tt.x, tt.y)
				if !ok {
					t.Fatalf("Подготовочный ход не удался")
				}
			}
			got := g.MakeMove(tt.x, tt.y)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("MakeMove(%d,%d) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestGame_NextPlayer(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		start int
		want  int
	}{
		{
			name:  "from 1 to 2",
			start: 1,
			want:  2,
		},
		{
			name:  "from 2 to 1",
			start: 2,
			want:  1,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &game.Game{
				Board:  nil,
				Player: tt.start,
			}
			g.NextPlayer()
			if g.Player != tt.want {
				t.Errorf("NextPlayer() Player = %d, want %d", g.Player, tt.want)
			}
		})
	}
}

func TestGame_Draw(t *testing.T) {
	tests := []struct {
		name  string // description of this test case
		board []string
		want  bool
	}{
		{
			name:  "true",
			board: []string{"X", "O", "X", "X", "O", "X", "O", "O"},
			want:  true,
		},
		{
			name:  "false",
			board: []string{"X", "O", "X", "X", "O", "X", "O", "#"},
			want:  false,
		},
		{
			name:  "nil-board",
			board: nil,
			want:  true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &game.Game{
				Board: tt.board,
			}
			got := g.Draw()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("Board: %v, Draw() = %v, want %v", tt.board, got, tt.want)
			}
		})
	}
}

func TestGame_Winner(t *testing.T) {
	tests := []struct {
		name  string
		board []string
		want  int
	}{
		{
			name:  "no winner empty board",
			board: []string{"#", "#", "#", "#", "#", "#", "#", "#", "#"},
			want:  0,
		},
		{
			name:  "X wins first row",
			board: []string{"X", "X", "X", "#", "#", "#", "#", "#", "#"},
			want:  1,
		},
		{
			name: "O wins second column",
			board: []string{
				"#", "O", "#",
				"#", "O", "#",
				"#", "O", "#",
			},
			want: 2,
		},
		{
			name: "X wins main diagonal",
			board: []string{
				"X", "#", "#",
				"#", "X", "#",
				"#", "#", "X",
			},
			want: 1,
		},
		{
			name: "O wins anti diagonal",
			board: []string{
				"#", "#", "O",
				"#", "O", "#",
				"O", "#", "#",
			},
			want: 2,
		},
		{
			name: "full board no winner",
			board: []string{
				"X", "O", "X",
				"X", "O", "O",
				"O", "X", "X",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &game.Game{
				Board: tt.board,
			}
			got := g.Winner()

			if got != tt.want {
				t.Errorf("Winner(%v) = %v, want %v", tt.board, got, tt.want)
			}
		})
	}
}
