package engine

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"errors"
	"fmt"
)

const (
	UP    string = "up"
	DOWN  string = "down"
	LEFT  string = "left"
	RIGHT string = "right"
	BOARD string = "b"
	NONE  string = "n"
)

var MOVES = [...]string{UP, DOWN, LEFT, RIGHT}
var OPPOSITE_MOVE = map[string]string{
	UP: DOWN,
	DOWN: UP,
	LEFT: RIGHT,
	RIGHT: LEFT,
}

func checkMove(g game.Game, coord game.Coordinate) (bool, error) {
	valid := true
	value := g.Board[coord.Y][coord.X]
	switch g.Board[coord.Y][coord.X] {
	case game.EMPTY:
		valid = true
	case game.FOOD:
		valid = true
	case game.WALL:
		valid = false
	default:
		if snakeValues, ok := g.ValueSnakeMap[value]; ok {
			size := len(snakeValues.Body)
			if coord.X == snakeValues.Body[size-1].X && coord.Y == snakeValues.Body[size-1].Y { // coord is a tail value and could be good or bad; evaluated during alpha-beta
				valid = true
			} else { // coord is a head or body value and means certain death
				valid = false
			}
		} else {
			return valid, errors.New("untracked board value")
		}
	}
	return valid, nil
}

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

func moveMeTo(g game.Game, newHeadCoord game.Coordinate) *node {
	fmt.Printf("MADE IT HERE")
	snake := g.ValueSnakeMap[game.ME]
	valid, err := checkMove(g, newHeadCoord)
	if valid && err == nil {
		g.Board[newHeadCoord.Y][newHeadCoord.X] = snake.Value
		if g.Board[newHeadCoord.Y][newHeadCoord.X] != game.FOOD {
			size := len(snake.Body)
			g.Board[snake.Body[size-1].Y][snake.Body[size-1].X] = game.EMPTY
			snake.Body = shiftBody(snake.Body, newHeadCoord)
			snake.Health -= 1
		} else {
			snake.Body = prependHead(snake.Body, newHeadCoord)
			snake.Health = 100
		}
		g.ValueSnakeMap[game.ME] = snake
		return &node{
			Game: g,
			Expanded: false,
		}
	} else {
		return nil
	}
}

func moveMe(g game.Game, direction string) *node {
	headCoord := g.ValueSnakeMap[game.ME].Body[0]
	switch direction {
	case UP:
		return moveMeTo(g, game.Coordinate{X:headCoord.X, Y:headCoord.Y-1})
	case DOWN:
		return moveMeTo(g, game.Coordinate{X:headCoord.X, Y:headCoord.Y+1})
	case LEFT:
		return moveMeTo(g, game.Coordinate{X:headCoord.X-1, Y:headCoord.Y})
	case RIGHT:
		return moveMeTo(g, game.Coordinate{X:headCoord.X+1, Y:headCoord.Y})
	default:
		return nil
	}
}

func generateMyGames(n node) []node {
	children  := make([]node, 3)
	successes := 0
	for _, move := range MOVES { // currently sequential, but may want to parallelize this after more testing
		if move == OPPOSITE_MOVE[n.Move] { continue }
		child := moveMe(n.Game, move)
		if child != nil {
			children[successes] = *child
			successes++
		}
	}
	children = children[:successes]
	return children
}

func generateBoardGames(n node) []node {
	return nil
}
