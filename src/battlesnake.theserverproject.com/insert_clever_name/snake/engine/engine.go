package engine

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/engine/evaluator"
	"battlesnake.theserverproject.com/insert_clever_name/snake/engine/expander"
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

func alphabeta(n *expander.Node, depth int, alpha float64, beta float64, maximizingPlayer bool) float64 {
	if depth == 0 || n.Terminal {
		return evaluator.Evaluate(n.Board)
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

func ComputeMove(g game.Game, deadline time.Duration) expander.Move {
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond) // process the move for x ms, leaving (500 - x) ms for the network
	defer cancel()
	root := expander.Node{
		Board:    game.CopyBoard(g.Board),
		Children: nil,
		Expanded: false,
	}
	latestMove := expander.UP // default move is some arbitrary direction for now
	depth := 2
	evaluated := make(chan int, 1)
	expanded := make(chan int, 1)
	var wg sync.WaitGroup
	evaluated <- 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context expired, therefore returning move", latestMove)
			return latestMove
		case <-evaluated:
			go func() {
				start := time.Now().UnixNano() / int64(time.Millisecond)
				wg.Add(1)
				expander.Expand(ctx, &root, depth, true, &wg)
				wg.Wait()
				end := time.Now().UnixNano() / int64(time.Millisecond)
				fmt.Printf("Expansion took %d milliseconds for depth %d\n", end - start, depth)
				expanded <- 1
			}()
		case <-expanded:
			go func() {
				start := time.Now().UnixNano() / int64(time.Millisecond)
				alphabeta(&root, depth, math.Inf(-1), math.Inf(1), true)
				end := time.Now().UnixNano() / int64(time.Millisecond)
				fmt.Printf("Evaluation took %d milliseconds for depth %d\n", end - start, depth)
				depth += 2
				evaluated <- 1
			}()
		}
	}
}
