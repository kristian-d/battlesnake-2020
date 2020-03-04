package engine

import (
	"math"
	"time"

	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
)

type node struct {
	Game     game.Game
	Children []node
	Expanded bool
	Move     string // this is the move type that generated this board from the previous board
}

func expandTree(n *node, depth int, maximizingPlayer bool) {
	if depth == 0 {
		return
	}
	if !n.Expanded {
		expandNode(n, maximizingPlayer)
	}
	for _, child := range n.Children {
		expandTree(&child, depth-1, !maximizingPlayer)
	}
}

func expandNode(n *node, maximizingPlayer bool) {
	n.Expanded = true
	if maximizingPlayer {
		n.Children = generateMyGames(*n)
	} else {
		n.Children = generateBoardGames(*n)
	}
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

func ComputeMove(g game.Game, deadline time.Duration) string {
	deadlineSignal := time.NewTimer(time.Millisecond * deadline).C // process the move for x ms, leaving (500 - x) ms for the network
	// some arbitrary depth for now. The initial depth should increase as the number of snakes decreases and size of snakes increases
	depth := 3
	root := node{
		Game:     g,
		Move:     NONE,
		Expanded: false,
	}

	latestMove := UP // default move is some arbitrary direction for now
	for {
		select {
		case <-deadlineSignal:
			return latestMove
		default:
			// TODO: compute new depth
			expandTree(&root, depth, true)
			//latestMove = alphabeta(root, depth, math.Inf(-1), math.Inf(1), true)
		}
	}
}
