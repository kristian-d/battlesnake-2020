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

func alphabeta(n *expander.Node, depth int, alpha float64, beta float64, maximizingPlayer bool) (float64, game.Move) {
	if depth == 0 || n.Terminal {
		return evaluator.Evaluate(n.Board), "" // with depths as multiples of two, this will always be a board move
	}
	if maximizingPlayer {
		value := math.Inf(-1) // negative infinity
		move := game.NONE
		for _, child := range n.Children {
			// generates my moves
			newValue, newMove := alphabeta(child, depth-1, alpha, beta, false)
			value = math.Max(value, newValue)
			if value == newValue {
				move = newMove
			}
			alpha = math.Max(alpha, value)
			if beta <= alpha {
				break
			}
		}
		return value, move
	} else {
		value := math.Inf(1) // positive infinity
		for _, child := range n.Children {
			// generates opponent's moves
			score, _ := alphabeta(child, depth-1, alpha, beta, true)
			value = math.Min(value, score)
			beta = math.Min(beta, value)
			if beta <= alpha {
				break
			}
		}
		return value, n.Board.MoveCoordinate.Move // return my move, which generated this board
	}
}

func ComputeMove(g game.Game, deadline time.Duration) game.Move {
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond) // process the move for x ms, leaving (500 - x) ms for the network
	defer cancel()
	absoluteDeadline := time.Now().UnixNano()/int64(time.Millisecond) + int64(deadline)
	root := expander.Node{
		Board:    game.CopyBoard(g.Board),
		Children: nil,
		Expanded: false,
	}
	latestMove := game.UP // default move is some arbitrary direction for now
	depth := g.PreviousMaxDepth - 2
	if depth < 2 {
		depth = 2
	}
	evaluated := make(chan int, 1)
	expanded := make(chan int, 1)
	lastExpansionElapsed := int64(0)
	var wg sync.WaitGroup
	evaluated <- 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context expired, therefore returning move", latestMove)
			g.PreviousMaxDepth = depth - 2
			return latestMove
		case <-evaluated:
			go func() {
				start := time.Now().UnixNano() / int64(time.Millisecond)
				wg.Add(1)
				expander.Expand(ctx, &root, depth, true, &wg)
				wg.Wait()
				end := time.Now().UnixNano() / int64(time.Millisecond)
				// fmt.Printf("Expansion took %d milliseconds for depth %d\n", end - start, depth) // for debugging
				lastExpansionElapsed = end - start
				expanded <- 1
			}()
		case <-expanded:
			go func() {
				// start := time.Now().UnixNano() / int64(time.Millisecond) // for debugging
				_, latestMove = alphabeta(&root, depth, math.Inf(-1), math.Inf(1), true)
				// end := time.Now().UnixNano() / int64(time.Millisecond) // for debugging
				// fmt.Printf("Evaluation took %d milliseconds for depth %d\n", end - start, depth) // for debugging
				depth += 2
				// this assumes that the next expansion will take longer than the previous expansion; if there isn't enough time, don't expand
				if absoluteDeadline - time.Now().UnixNano()/int64(time.Millisecond) > lastExpansionElapsed {
					evaluated <- 1
				}
			}()
		}
	}
}
