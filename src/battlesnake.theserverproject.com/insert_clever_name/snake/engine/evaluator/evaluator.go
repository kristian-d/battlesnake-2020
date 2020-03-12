package evaluator

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"math/rand"
)

func Evaluate(b game.Board) float64 {
	if _, ok := b.Snakes[game.ME]; !ok {
		return 0 // if we are dead, return minimum evaluation
	} else if len(b.Snakes) == 1 {
		return 100 // if we are the only snake alive, return maximum evaluation
	}
	return float64(rand.Intn(100))
}
