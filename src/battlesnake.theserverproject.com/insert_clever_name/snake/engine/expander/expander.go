package expander

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"context"
	"math"
	"sync"
)

type Node struct {
	Game     game.Game
	Children []*Node
	Expanded bool
	Move     Move // this is the move type that generated this board from the previous board
	Terminal bool
}

func gameBranchesBySnakeMove(g game.Game, snakeValue game.BoardValue) <-chan game.Game {
	// the buffer prevents any of the go routines from hanging if the receiver stops listening
	c := make(chan game.Game, 3)
	head := g.ValueSnakeMap[snakeValue].Body[0]
	newHeadCoords := [...]game.Coordinate{
		{X:head.X, Y:head.Y - 1}, // UP
		{X:head.X, Y:head.Y + 1}, // DOWN
		{X:head.X - 1, Y:head.Y}, // LEFT
		{X:head.X + 1, Y:head.Y}, // RIGHT
	}

	var wg sync.WaitGroup
	wg.Add(4)
	output := func(coord game.Coordinate) {
		defer wg.Done()
		newGame := game.CopyGame(g)
		moveSnake(newGame, snakeValue, coord)
		c <- newGame
	}

	successful := 0
	for _, coord := range newHeadCoords {
		if prelimaryCheck(g, snakeValue, coord) {
			successful++
			go output(coord)
		} else {
			wg.Done()
		}
	}

	go func() {
		wg.Wait()
		// if the snake could not move anywhere, its death is the only path
		if successful == 0 {
			newGame := game.CopyGame(g)
			killSnake(newGame, snakeValue)
			c <- newGame
		}
		close(c)
	}()

	return c
}

func gameBranches(g game.Game) <-chan game.Game {
	valueSnakeMap := g.ValueSnakeMap
	maxSize := 0
	var largestSnakeValue game.BoardValue
	for value, snake := range valueSnakeMap {
		if !snake.Moved && len(snake.Body) > maxSize {
			maxSize = len(snake.Body)
			largestSnakeValue = value
		}
	}
	return gameBranchesBySnakeMove(g, largestSnakeValue)
}

func nextGameStates(ctx context.Context, g game.Game, maximizingPlayer bool) <-chan game.Game {
	if maximizingPlayer {
		return gameBranchesBySnakeMove(g, game.ME)
	}
	// buffer channels to the maximum possible number of outputs so that there are no blocks
	maxOutputs := int64(math.Pow(3, float64(len(g.ValueSnakeMap)-1)))
	in := make(chan game.Game, maxOutputs)
	out := make(chan game.Game, maxOutputs)
	var wg sync.WaitGroup

	pipe := func(branches <-chan game.Game) {
		defer wg.Done()
		for branch := range branches {
			if turnComplete(branch) {
				select {
				case out <- branch:
				case <-ctx.Done():
					return
				}
			} else {
				wg.Add(1)
				select {
				case in <- branch:
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
		}
	}

	go func() {
		for branch := range in {
			go pipe(gameBranches(branch))
		}
	}()

	wg.Add(1)
	in <- g
	go func() {
		wg.Wait()
		close(in)
		close(out)
	}()
	return out
}

func Expand(ctx context.Context, n *Node, depth int, maximizingPlayer bool, wg *sync.WaitGroup) {
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
				childRef := &Node{
					Game: childGame,
					Children: nil,
					Expanded: false,
					Terminal: terminalNode,
				}
				n.Children = append(n.Children, childRef)
				wg.Add(1)
				go Expand(ctx, childRef, depth-1, !maximizingPlayer, wg)
			}
		}
		n.Expanded = true
		n.Game.Board = nil
		n.Game.ValueSnakeMap = nil
	} else {
		for _, child := range n.Children {
			wg.Add(1)
			go Expand(ctx, child, depth-1, !maximizingPlayer, wg)
		}
	}
}
