package expander

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"context"
	"math"
	"sync"
)

type Node struct {
	Board    game.Board
	Children []*Node
	Expanded bool
	Move     game.Move // this is the move type that generated this board from the previous board
	Terminal bool
}

func boardBranchesBySnakeMove(b game.Board, snakeValue game.GridValue) <-chan game.Board {
	// the buffer prevents any of the go routines from hanging if the receiver stops listening
	c := make(chan game.Board, 3)
	head := b.Snakes[snakeValue].Body[0]
	moveCoords := [...]game.MoveCoordinate{
		{game.UP, game.Coordinate{X:head.X, Y:head.Y - 1}},
		{game.DOWN, game.Coordinate{X:head.X, Y:head.Y + 1}},
		{game.LEFT, game.Coordinate{X:head.X - 1, Y:head.Y}},
		{game.RIGHT, game.Coordinate{X:head.X + 1, Y:head.Y}},
	}

	var wg sync.WaitGroup
	wg.Add(4)
	output := func(moveCoord game.MoveCoordinate) {
		defer wg.Done()
		newBoard := game.CopyBoard(b)
		moveSnake(newBoard, snakeValue, moveCoord.Coordinate)
		newBoard.MoveCoordinate = moveCoord
		c <- newBoard
	}

	successful := 0
	for _, moveCoord := range moveCoords {
		if prelimaryCheck(b, snakeValue, moveCoord.Coordinate) {
			successful++
			go output(moveCoord)
		} else {
			wg.Done()
		}
	}

	go func() {
		wg.Wait()
		// if the snake could not move anywhere, its death is the only path
		if successful == 0 {
			newBoard := game.CopyBoard(b)
			killSnake(newBoard, snakeValue)
			c <- newBoard
		}
		close(c)
	}()

	return c
}

func boardBranches(b game.Board) <-chan game.Board {
	valueSnakeMap := b.Snakes
	maxSize := 0
	var largestSnakeValue game.GridValue
	for value, snake := range valueSnakeMap {
		if !snake.Moved && len(snake.Body) > maxSize {
			maxSize = len(snake.Body)
			largestSnakeValue = value
		}
	}
	return boardBranchesBySnakeMove(b, largestSnakeValue)
}

func nextBoardStates(ctx context.Context, b game.Board, maximizingPlayer bool) <-chan game.Board {
	if maximizingPlayer {
		return boardBranchesBySnakeMove(b, game.ME)
	}
	// buffer channels to the maximum possible number of outputs so that there are no blocks
	maxOutputs := int64(math.Pow(3, float64(len(b.Snakes)-1)))
	in := make(chan game.Board, maxOutputs)
	out := make(chan game.Board, maxOutputs)
	var wg sync.WaitGroup

	pipe := func(branches <-chan game.Board) {
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
			go pipe(boardBranches(branch))
		}
	}()

	wg.Add(1)
	in <- b
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
			resetTurn(n.Board)
		}
		for childBoard := range nextBoardStates(ctx, n.Board, maximizingPlayer) {
			select {
			case <-ctx.Done():
				return
			default:
				terminalNode := false
				// if we are dead or we are the only snake left, path is terminated
				if _, ok := childBoard.Snakes[game.ME]; !ok || len(childBoard.Snakes) == 1 {
					terminalNode = true
				}
				childRef := &Node{
					Board:    childBoard,
					Children: nil,
					Expanded: false,
					Terminal: terminalNode,
				}
				n.Children = append(n.Children, childRef)
				wg.Add(1)
				go Expand(ctx, childRef, depth-1, !maximizingPlayer, wg)
			}
		}
		n.Expanded     = true
		n.Board.Grid   = nil
		n.Board.Snakes = nil
	} else {
		for _, child := range n.Children {
			wg.Add(1)
			go Expand(ctx, child, depth-1, !maximizingPlayer, wg)
		}
	}
}
