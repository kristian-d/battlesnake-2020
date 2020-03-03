package engine

import (
	"math"
	"time"

	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
)

const (
	UP    string = "up"
	DOWN  string = "down"
	LEFT  string = "left"
	RIGHT string = "right"
)

type node struct {
	Board    [][]int
	Children []node
	Expanded bool
}

func expandTree(root *node, depth int, expanded chan<- int) {
	expanded <- 0
	return
}

func evaluate(n node) float64 {
	return 0
}

func alphabeta(n node, depth int, alpha float64, beta float64, maximizingPlayer bool) float64 {
	if depth == 0 || len(n.Children) == 0 {
		return evaluate(n)
	}
	if maximizingPlayer {
		value := math.Inf(-1) // negative infinity
		for _, child := range n.Children {
			value = math.Max(value, alphabeta(child, depth-1, alpha, beta, false))
			alpha = math.Max(alpha, value)
			if beta <= alpha {
				break
			}
		}
		return value
	} else {
		value := math.Inf(1) // positive infinity
		for _, child := range n.Children {
			value = math.Min(value, alphabeta(child, depth-1, alpha, beta, true))
			beta = math.Min(beta, value)
			if beta <= alpha {
				break
			}
		}
		return value
	}
}

func ComputeMove(game *game.Game, deadline time.Duration) string {
	deadlineSignal := time.NewTimer(time.Millisecond * deadline).C // process the move for x ms, leaving (500 - x) ms for the network
	// some arbitrary depth for now. The initial depth should increase as the number of snakes decreases and size of snakes increases
	depth := 3
	root := node{
		Board:    game.Board,
		Children: make([]node, 0),
	}
	expanded := make(chan int, 1)
	/*computed := make(chan string, 1)
	expandTree(&root, depth, expanded)
	for {
		select {
		case <-expanded:
			// TODO: compute new depth
			go alphabeta(root, depth, math.Inf(-1), math.Inf(1), true)
			go expandTree(&root, depth, expanded)
		case move := <-computed:
			moveChan <- move
		case <-quitChan:
			return
		}
	}*/

	// latestMove is flagged as unneeded, and thus the
	// project will not compile. Commenting it for now.
	// TODO: Use or remove the below line of code.
	//latestMove := UP // default move is some arbitrary direction for now
	for {
		select {
		case <-deadlineSignal:
			return UP
		default:
			// TODO: compute new depth
			expandTree(&root, depth, expanded)
			//latestMove := alphabeta(root, depth, math.Inf(-1), math.Inf(1), true)
		}
	}
}