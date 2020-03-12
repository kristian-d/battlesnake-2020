package expander

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
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
	height := len(g.Grid)
	width := len(g.Grid[0])
	return !(coord.X < width && coord.X >= 0 && coord.Y < height && coord.Y >= 0)
}

func killSnake(g game.Game, snakeValue game.GridValue) {
	grid := g.Grid
	snake := g.ValueSnakeMap[snakeValue]
	for _, bodyPart := range snake.Body {
		grid[bodyPart.Y][bodyPart.X] = game.EMPTY
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

func resetTurn(g game.Game) {
	for value, snake := range g.ValueSnakeMap {
		snake.Moved = false
		g.ValueSnakeMap[value] = snake
	}
}

func prelimaryCheck(g game.Game, snakeValue game.GridValue, coord game.Coordinate) bool {
	if outOfBounds(g, coord) {
		return false // moving off of the grid, therefore guaranteed death
	}
	value := g.Grid[coord.Y][coord.X]
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
