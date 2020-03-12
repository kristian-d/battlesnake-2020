package expander

import "battlesnake.theserverproject.com/insert_clever_name/snake/game"

func moveSnake(b game.Board, snakeValue game.GridValue, coord game.Coordinate) {
	grid := b.Grid
	value := grid[coord.Y][coord.X]
	snake := b.Snakes[snakeValue]
	size := len(snake.Body)
	snake.Moved = true
	if value != game.FOOD {
		if otherSnake, ok := b.Snakes[value]; ok && coord.X == otherSnake.Body[0].X && coord.Y == otherSnake.Body[0].Y { // moving onto a head value
			killSnake(b, value) // this will be me if all other snakes are handled in order of decreasing size
		}
		snake.Health -= 1
		// if tail location is still tail value, then set it to empty, else another snake's head has already moved there
		if grid[snake.Body[size-1].Y][snake.Body[size-1].X] == snakeValue {
			grid[snake.Body[size-1].Y][snake.Body[size-1].X] = game.EMPTY
		}
		snake.Body = shiftBody(snake.Body, coord)
	} else {
		snake.Health = 100
		snake.Body = prependHead(snake.Body, coord)
		// if grown and tail value is not own value, then another snake's head has moved onto tail and must die
		if tailValue := grid[snake.Body[size-1].Y][snake.Body[size-1].X]; tailValue != snakeValue {
			killSnake(b, tailValue)
			grid[snake.Body[size-1].Y][snake.Body[size-1].X] = snakeValue
		}
	}
	grid[coord.Y][coord.X] = snake.Value
	b.Snakes[snakeValue] = snake
}
