package expander

import "battlesnake.theserverproject.com/insert_clever_name/snake/game"

type Move string
const (
	UP    Move = "up"
	DOWN  Move = "down"
	LEFT  Move = "left"
	RIGHT Move = "right"
)

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
