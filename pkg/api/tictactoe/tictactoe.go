package tictactoe

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type TicTacToe struct {
	Board       map[string]string
	CurrentTurn string
	IsGameOver  bool
	Mutex       sync.Mutex
}

func New() *TicTacToe {
	return &TicTacToe{
		Board: map[string]string{
			"1": "1", "2": "2", "3": "3",
			"4": "4", "5": "5", "6": "6",
			"7": "7", "8": "8", "9": "9",
		},
		CurrentTurn: "X",
		IsGameOver:  false,
		Mutex:       sync.Mutex{},
	}
}

var emojiMap = map[string]string{
	"X": "❌",
	"O": "⭕",
	"1": "1️⃣",
	"2": "2️⃣",
	"3": "3️⃣",
	"4": "4️⃣",
	"5": "5️⃣",
	"6": "6️⃣",
	"7": "7️⃣",
	"8": "8️⃣",
	"9": "9️⃣",
}

func (t *TicTacToe) PrintBoard() string {
	return fmt.Sprintf(`
     %s | %s | %s
    ---------------
     %s | %s | %s
    ---------------
     %s | %s | %s
    `, t.getEmoji("1"), t.getEmoji("2"), t.getEmoji("3"),
		t.getEmoji("4"), t.getEmoji("5"), t.getEmoji("6"),
		t.getEmoji("7"), t.getEmoji("8"), t.getEmoji("9"))
}

func (t *TicTacToe) getEmoji(position string) string {
	if val, ok := t.Board[position]; ok {
		if emoji, emojiExists := emojiMap[val]; emojiExists {
			return emoji
		}
	}
	return position
}

func (t *TicTacToe) MakeMove(position string) bool {
	if t.Board[position] == position {
		t.Board[position] = t.CurrentTurn
		return true
	}
	return false
}

func (t *TicTacToe) CheckWin() bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if t.Board[strconv.Itoa(i*3+1)] == t.CurrentTurn &&
			t.Board[strconv.Itoa(i*3+2)] == t.CurrentTurn &&
			t.Board[strconv.Itoa(i*3+3)] == t.CurrentTurn {
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if t.Board[strconv.Itoa(i+1)] == t.CurrentTurn &&
			t.Board[strconv.Itoa(i+4)] == t.CurrentTurn &&
			t.Board[strconv.Itoa(i+7)] == t.CurrentTurn {
			return true
		}
	}

	// Check diagonals
	if t.Board["1"] == t.CurrentTurn && t.Board["5"] == t.CurrentTurn && t.Board["9"] == t.CurrentTurn {
		return true
	}

	if t.Board["3"] == t.CurrentTurn && t.Board["5"] == t.CurrentTurn && t.Board["7"] == t.CurrentTurn {
		return true
	}

	return false
}

func (t *TicTacToe) CheckDraw() bool {
	for _, value := range t.Board {
		if value != "X" && value != "O" {
			return false
		}
	}
	return true
}

func (t *TicTacToe) SwitchTurn() {
	if t.CurrentTurn == "X" {
		t.CurrentTurn = "O"
	} else {
		t.CurrentTurn = "X"
	}
}

func (t *TicTacToe) ApplyBotMove(move string) {
	var updatedBoard map[string]string
	err := json.Unmarshal([]byte(move), &updatedBoard)
	if err != nil {
		fmt.Println("Error decoding bot's move:", err)
		return
	}

	t.Board = updatedBoard
}

func (t *TicTacToe) GetStateString() string {
	state, err := json.Marshal(t.Board)
	if err != nil {
		fmt.Println("Error encoding game state:", err)
		return ""
	}
	return string(state)
}

func (t *TicTacToe) GetRandBotMove() string {
	rand.Seed(time.Now().UnixNano())
	validMoves := t.GetValidMoves()
	if len(validMoves) == 0 {
		return ""
	}
	return validMoves[rand.Intn(len(validMoves))]
}

func (t *TicTacToe) GetValidMoves() []string {
	var validMoves []string
	for position := 1; position <= 9; position++ {
		if t.Board[strconv.Itoa(position)] == strconv.Itoa(position) {
			validMoves = append(validMoves, strconv.Itoa(position))
		}
	}
	return validMoves
}
