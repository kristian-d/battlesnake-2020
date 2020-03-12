package evaluator

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"math/rand"
)

func Evaluate(b game.Board) float64 {
	return float64(rand.Intn(100))
}
