package engine

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
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

func expandTree(ctx context.Context, n *node, depth int, maximizingPlayer bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth == 0 || n.Terminal {
		return
	}
	if !n.Expanded {
		if maximizingPlayer {
			resetTurn(n.Game)
		}
		for childGame := range nextGameStates(ctx, n.Game, maximizingPlayer) {
			select {
			case <-ctx.Done():
				return
			default:
				terminalNode := false
				// if we are dead or we are the only snake left, path is terminated
				if _, ok := childGame.ValueSnakeMap[game.ME]; !ok || len(childGame.ValueSnakeMap) == 1 {
					terminalNode = true
				}
				childRef := &node{
					Game: childGame,
					Children: nil,
					Expanded: false,
					Terminal: terminalNode,
				}
				n.Children = append(n.Children, childRef)
				wg.Add(1)
				go expandTree(ctx, childRef, depth-1, !maximizingPlayer, wg)
			}
		}
		n.Expanded = true
		n.Game.Board = nil
		n.Game.ValueSnakeMap = nil
	} else {
		for _, child := range n.Children {
			wg.Add(1)
			go expandTree(ctx, child, depth-1, !maximizingPlayer, wg)
		}
	}
}

func evaluate(g game.Game) float64 {
	return float64(rand.Intn(100))
}

func alphabeta(n *node, depth int, alpha float64, beta float64, maximizingPlayer bool) float64 {
	if depth == 0 || n.Terminal {
		return evaluate(n.Game)
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
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond) // process the move for x ms, leaving (500 - x) ms for the network
	defer cancel()
	root := node{
		Game:     g,
		Children: nil,
		Expanded: false,
	}
	latestMove := UP // default move is some arbitrary direction for now
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
				expandTree(ctx, &root, depth, true, &wg)
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
