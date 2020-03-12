package evaluator

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"math/rand"
)

func Evaluate(g game.Game) float64 {
	return float64(rand.Intn(100))
}
