package engine

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
)

type node struct {
	Game     game.Game
	Children []*node
	Expanded bool
	Move     Move // this is the move type that generated this board from the previous board
	Terminal bool
}

func expandTree(done <-chan int, n *node, depth int, maximizingPlayer bool) {
	if depth == 0 || n.Terminal {
		return
	}
	if !n.Expanded {
		if maximizingPlayer {
			resetTurn(n.Game)
		}
		for child := range expandNode(done, *n, maximizingPlayer) {
			n.Children = append(n.Children, &child)
			expandTree(done, &child, depth-1, !maximizingPlayer)
		}
		n.Expanded = true
		n.Game.Board = nil
		n.Game.ValueSnakeMap = nil
	} else {
		for _, child := range n.Children {
			expandTree(done, child, depth-1, !maximizingPlayer)
		}
	}
}

func expandNode(done <-chan int, n node, maximizingPlayer bool) <-chan node {
	out := make(chan node, int64(math.Pow(3, float64(len(n.Game.ValueSnakeMap)-1))))
	go func() {
		defer close(out)
		for child := range nextGameStates(done, n.Game, maximizingPlayer) {
			terminalNode := false
			// if we are dead or we are the only snake left, path is terminated
			if _, ok := child.ValueSnakeMap[game.ME]; !ok || len(child.ValueSnakeMap) == 1 {
				terminalNode = true
			}
			out <- node{
				Game: child,
				Children: nil,
				Expanded: false,
				Terminal: terminalNode,
			}
		}
	}()
	return out
}

func evaluate(n node) float64 {
	return float64(rand.Intn(100))
}

func alphabeta(n *node, depth int, alpha float64, beta float64, maximizingPlayer bool) float64 {
	if depth == 0 || n.Terminal {
		return evaluate(*n)
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
	deadlineSignal := time.NewTimer(time.Millisecond * deadline).C // process the move for x ms, leaving (500 - x) ms for the network
	root := node{
		Game:     g,
		Children: nil,
		Expanded: false,
	}
	latestMove := UP // default move is some arbitrary direction for now
	done := make(chan int)
	depth := 2
	for {
		select {
		case <-deadlineSignal:
			done <- 1
			return latestMove
		default:
			start := time.Now().UnixNano() / int64(time.Millisecond)
			expandTree(done, &root, depth, true)
			score := alphabeta(&root, depth, math.Inf(-1), math.Inf(1), true)
			end := time.Now().UnixNano() / int64(time.Millisecond)
			fmt.Printf("Evaluation took %d milliseconds for depth %d, resulting in a score of %f\n", end - start, depth, score)
			depth += 2
		}
	}
}
