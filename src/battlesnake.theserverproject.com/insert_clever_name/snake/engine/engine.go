package engine

import (
	"fmt"
	"math"
	"time"

	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
)

type node struct {
	Game     game.Game
	Children []node
	Expanded bool
	Move     Move // this is the move type that generated this board from the previous board
}

func expandTree(done <-chan int, n *node, depth int, maximizingPlayer bool) {
	if depth == 0 {
		return
	}
	if !n.Expanded {
		resetTurn(n.Game)
		children := expandNode(done, *n, maximizingPlayer)
		for child := range children {
			n.Children = append(n.Children, child)
			expandTree(done, &child, depth-1, !maximizingPlayer)
		}
		n.Expanded = true
		n.Game.Board = nil
		n.Game.ValueSnakeMap = nil
	} else {
		for _, child := range n.Children {
			expandTree(done, &child, depth-1, !maximizingPlayer)
		}
	}
}

func expandNode(done <-chan int, n node, maximizingPlayer bool) <-chan node {
	out := make(chan node)
	go func() {
		defer close(out)
		c := nextGameStates(done, n.Game, maximizingPlayer)
		for child := range c {
			out <- node{
				Game: child,
				Children: nil,
				Expanded: false,
			}
		}
	}()
	return out
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

func ComputeMove(g game.Game, deadline time.Duration) Move {
	//deadlineSignal := time.NewTimer(time.Millisecond * deadline).C // process the move for x ms, leaving (500 - x) ms for the network
	// some arbitrary depth for now. The initial depth should increase as the number of snakes decreases and size of snakes increases
	depth := 6
	root := node{
		Game:     g,
		Children: nil,
		Expanded: false,
	}

	latestMove := UP // default move is some arbitrary direction for now
	done := make(chan int)
	defer close(done)
	expandTree(done, &root, depth, true)
	fmt.Println("Done")
	return latestMove
}
