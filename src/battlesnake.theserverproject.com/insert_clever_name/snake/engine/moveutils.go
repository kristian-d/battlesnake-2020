package engine

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"sync"
)

type Move string
const (
	UP    Move = "up"
	DOWN  Move = "down"
	LEFT  Move = "left"
	RIGHT Move = "right"
)

func prependHead(body []game.Coordinate, head game.Coordinate) []game.Coordinate { // theoretically faster than using straight append(), should be tested
	body = append(body, head)
	copy(body[1:], body)
	body[0] = head
	return body
}

func shiftBody(body []game.Coordinate, head game.Coordinate) []game.Coordinate {
	for i := len(body) - 1; i >= 1; i-- {
		body[i] = body[i - 1]
	}
	body[0] = head
	return body
}

func outOfBounds(g game.Game, coord game.Coordinate) bool {
	height := len(g.Board)
	width := len(g.Board[0])
	return !(coord.X < width && coord.X >= 0 && coord.Y < height && coord.Y >= 0)
}

func killSnake(g game.Game, snakeValue game.BoardValue) {
	board := g.Board
	snake := g.ValueSnakeMap[snakeValue]
	for _, bodyPart := range snake.Body {
		board[bodyPart.Y][bodyPart.X] = game.EMPTY
	}
	delete(g.ValueSnakeMap, snakeValue)
}

func turnComplete(g game.Game) bool {
	for _, snake := range g.ValueSnakeMap {
		if !snake.Moved && snake.Value != game.ME {
			return false
		}
	}
	return true
}

func myTurnComplete(g game.Game) bool {
	return g.ValueSnakeMap[game.ME].Moved
}

func resetTurn(g game.Game) game.Game {
	for value, snake := range g.ValueSnakeMap {
		snake.Moved = false
		g.ValueSnakeMap[value] = snake
	}
	return g
}

func prelimaryCheck(g game.Game, snakeValue game.BoardValue, coord game.Coordinate) bool {
	if outOfBounds(g, coord) {
		return false // moving off of the board, therefore guaranteed death
	}
	value := g.Board[coord.Y][coord.X]
	if value == game.FOOD {
		return true // moving into a location with food, therefore death is not guaranteed
	}
	if g.ValueSnakeMap[snakeValue].Health == 1 {
		return false // starvation is next turn and location does not contain food, therefore guaranteed death
	}
	if value == game.EMPTY {
		return true // moving into empty location, therefore death is not guaranteed
	}
	otherSnake := g.ValueSnakeMap[value]
	if otherSnake.Moved {
		if coord.X == otherSnake.Body[0].X && coord.Y == otherSnake.Body[0].Y && len(g.ValueSnakeMap[snakeValue].Body) > len(otherSnake.Body) {
			return true // moving onto a head value and has size advantage, therefore death is not guaranteed
		} else {
			return false // moving onto a body value, tail value, or head value without size advantage, therefore guaranteed death
		}
	} else {
		if coord.X == otherSnake.Body[len(otherSnake.Body)-1].X && coord.Y == otherSnake.Body[len(otherSnake.Body)-1].Y {
			return true // moving onto a tail value, therefore death is not guaranteed
		} else {
			return false // moving onto a body value, therefore guaranteed death
		}
	}
}

func moveSnake(g game.Game, snakeValue game.BoardValue, coord game.Coordinate) {
	board := g.Board
	value := board[coord.Y][coord.X]
	snake := g.ValueSnakeMap[snakeValue]
	size := len(snake.Body)
	snake.Moved = true
	if value != game.FOOD {
		if otherSnake, ok := g.ValueSnakeMap[value]; ok && coord.X == otherSnake.Body[0].X && coord.Y == otherSnake.Body[0].Y { // moving onto a head value
			killSnake(g, value) // this will be me if all other snakes are handled in order of decreasing size
		}
		snake.Health -= 1
		// if tail location is still tail value, then set it to empty, else another snake's head has already moved there
		if board[snake.Body[size-1].Y][snake.Body[size-1].X] == snakeValue {
			board[snake.Body[size-1].Y][snake.Body[size-1].X] = game.EMPTY
		}
		snake.Body = shiftBody(snake.Body, coord)
	} else {
		snake.Health = 100
		snake.Body = prependHead(snake.Body, coord)
		// if grown and tail value is not own value, then another snake's head has moved onto tail and must die
		if tailValue := board[snake.Body[size-1].Y][snake.Body[size-1].X]; tailValue != snakeValue {
			killSnake(g, tailValue)
			board[snake.Body[size-1].Y][snake.Body[size-1].X] = snakeValue
		}
	}
	board[coord.Y][coord.X] = snake.Value
	g.ValueSnakeMap[snakeValue] = snake
}

func gameBranchesBySnakeMove(g game.Game, snakeValue game.BoardValue) <-chan game.Game {
	// the buffer prevents any of the four go routines from hanging if the receiver stops listening
	c := make(chan game.Game, 4)
	head := g.ValueSnakeMap[snakeValue].Body[0]
	newHeadCoords := [...]game.Coordinate{
		{X:head.X, Y:head.Y - 1}, // UP
		{X:head.X, Y:head.Y + 1}, // DOWN
		{X:head.X - 1, Y:head.Y}, // LEFT
		{X:head.X + 1, Y:head.Y}, // RIGHT
	}

	var wg sync.WaitGroup
	wg.Add(4)
	for _, coord := range newHeadCoords {
		go func() {
			defer wg.Done()
			if prelimaryCheck(g, snakeValue, coord) {
				newGame := game.CopyGame(g)
				moveSnake(newGame, snakeValue, coord)
				c <- newGame
			}
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	return c
}

func myGameBranches(g game.Game) <-chan game.Game {
	return gameBranchesBySnakeMove(g, game.ME)
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

func nextGameStates(done <-chan int, g game.Game, maximizingPlayer bool) <-chan game.Game {
	sink := make(chan game.Game)
	feedback := make(chan game.Game)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer close(sink)
		for branch := range feedback {
			go func() {
				defer wg.Done()
				var branches <-chan game.Game
				if maximizingPlayer {
					branches = myGameBranches(branch)
				} else {
					branches = gameBranches(branch)
				}
				for newBranch := range branches {
					if (maximizingPlayer && myTurnComplete(newBranch)) || (!maximizingPlayer && turnComplete(newBranch)) {
						select {
						case sink <- newBranch:
						case <-done:
							return
						}
					} else {
						wg.Add(1)
						select {
						case feedback <- newBranch:
						case <-done:
							return
						}
					}
				}
			}()
		}

	}()

	go func() {
		wg.Wait()
		close(feedback)
	}()

	feedback <- g
	return sink
}
