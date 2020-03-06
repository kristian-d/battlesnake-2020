package engine

import (
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
	"errors"
)

type Move string
const (
	UP    Move = "up"
	DOWN  Move = "down"
	LEFT  Move = "left"
	RIGHT Move = "right"
)
var moves = [...]Move{UP, DOWN, LEFT, RIGHT}
var oppositeMove = map[Move]Move{
	UP:    DOWN,
	DOWN:  UP,
	LEFT:  RIGHT,
	RIGHT: LEFT,
}

type Event int
const (
	NONE Event = iota
	EAT
	DEATH
	RISK
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

func checkMyMove(g game.Game, coord game.Coordinate) (bool, error) {
	value := g.Board[coord.Y][coord.X]
	switch g.Board[coord.Y][coord.X] {
	case game.EMPTY:
		return true, nil
	case game.FOOD:
		return true, nil
	case game.WALL:
		return false, nil
	default:
		if snakeValues, ok := g.ValueSnakeMap[value]; ok {
			size := len(snakeValues.Body)
			if coord.X == snakeValues.Body[size-1].X && coord.Y == snakeValues.Body[size-1].Y { // coord is a tail value and could be good or bad; evaluated during alpha-beta
				return true, nil
			} else { // coord is a head or body value and means certain death
				return false, nil
			}
		} else {
			return false, errors.New("untracked board value")
		}
	}
}

func moveMeTo(g game.Game, newHeadCoord game.Coordinate) *node {
	snake := g.ValueSnakeMap[game.ME]
	valid, err := checkMyMove(g, newHeadCoord)
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

func checkEvent(g game.Game, coord game.Coordinate) (Event, error) {
	value := g.Board[coord.Y][coord.X]
	switch value {
	case game.EMPTY:
		return NONE, nil
	case game.FOOD:
		return EAT, nil
	case game.WALL:
		return DEATH, nil
	default:
		if snakeValues, ok := g.ValueSnakeMap[value]; ok {
			size := len(snakeValues.Body)
			if coord.X == snakeValues.Body[size-1].X && coord.Y == snakeValues.Body[size-1].Y { // coord is a tail value and is risky; evaluated during alpha-beta
				return RISK, nil
			} else { // coord is a head or body value and means certain death
				return DEATH, nil
			}
		} else {
			return NONE, errors.New("untracked board value")
		}
	}
}

func moveSnakeTo(g game.Game, snakeValue game.BoardValue, newHeadCoord game.Coordinate) *node {
	event, err := checkEvent(g, newHeadCoord)
	if err != nil { // TODO: perhaps do some other handling here
		return nil
	}
	snake := g.ValueSnakeMap[snakeValue]
	switch event {
	case NONE:
		g.Board[newHeadCoord.Y][newHeadCoord.X] = snake.Value
		size := len(snake.Body)
		g.Board[snake.Body[size-1].Y][snake.Body[size-1].X] = game.EMPTY
		snake.Body = shiftBody(snake.Body, newHeadCoord)
		snake.Health -= 1
	case EAT:
		g.Board[newHeadCoord.Y][newHeadCoord.X] = snake.Value
		snake.Body = prependHead(snake.Body, newHeadCoord)
		snake.Health = 100
	case DEATH:
		if snakeValue == game.ME { // death isn't an option! Chaaarge!
			return nil
		}
		// TODO: handle death of snake
	case RISK:
	}
	g.ValueSnakeMap[snakeValue] = snake
	return &node{
		Game: g,
		Expanded: false,
	}
}

func moveMe(g game.Game, direction Move) *node {
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
	for _, move := range moves { // currently sequential, but may want to parallelize this after more testing
		if move == oppositeMove[n.Move] { continue }
		child := moveMe(n.Game, move)
		if child != nil {
			children[successes] = *child
			successes++
		}
	}
	children = children[:successes]
	return children
}

/*func generatePossibleMoves(g game.Game, snake game.Snake) []game.Coordinate {
	head := snake.Body[0]
	tail := snake.Body[len(snake.Body) - 1]
	newHeadCoords := [...]game.Coordinate{
		{X:head.X, Y:head.Y - 1}, // UP
		{X:head.X, Y:head.Y + 1}, // DOWN
		{X:head.X - 1, Y:head.Y}, // LEFT
		{X:head.X + 1, Y:head.Y}, // RIGHT
	}
	moveCoords := make([]game.Coordinate, 3) // three is the maximum possible open spaces around the head
	successes := 0
	for _, newHead := range newHeadCoords {
		boardValue := g.Board[newHead.Y][newHead.X]
		// assume opponent snake won't go back into body or into a wall voluntarily unless forced to
		if (boardValue == snake.Value && !(newHead.X == tail.X || newHead.Y == tail.Y)) ||
			boardValue == game.WALL {
			continue
		}
		moveCoords[successes] = newHead
		successes++
	}
	if len(moveCoords) == 0 { // snake is forced to kill itself
		moveCoords[0] = newHeadCoords[0]
		successes = 1
	}
	return moveCoords[:successes]
}*/

type Moveset map[game.Coordinate][]game.BoardValue

func combineMovesets(m1, m2 Moveset) Moveset {
	for coord, snakeValues := range m1 {
		coordsCopy := copy(make([]game.BoardValue, len(snakeValues)), snakeValues)
		if _, ok := m2[coord]; ok {
			m2[coord] = append(m2[coord], coordsCopy...)
		}
	}
}

func generatePossibleMoves(g game.Game, snake game.Snake) []Moveset {
	head := snake.Body[0]
	tail := snake.Body[len(snake.Body) - 1]
	newHeadCoords := [...]game.Coordinate{
		{X:head.X, Y:head.Y - 1}, // UP
		{X:head.X, Y:head.Y + 1}, // DOWN
		{X:head.X - 1, Y:head.Y}, // LEFT
		{X:head.X + 1, Y:head.Y}, // RIGHT
	}
	movesets := make([]Moveset, 3)
	successes := 0
	for _, newHead := range newHeadCoords {
		boardValue := g.Board[newHead.Y][newHead.X]
		// assume opponent snake won't go back into body or into a wall voluntarily unless forced to
		if (boardValue == snake.Value && !(newHead.X == tail.X || newHead.Y == tail.Y)) ||
			boardValue == game.WALL {
			continue
		}
		movesets[successes][newHead] = []game.BoardValue{snake.Value}
		successes++
	}
	if len(movesets) == 0 { // snake is forced to kill itself
		movesets[0][newHeadCoords[0]] = []game.BoardValue{snake.Value}
		successes = 1
	}
	return movesets[:successes]
}

func generateMovesets(g game.Game, valueSnakeMap game.SnakeByValue) <-chan Moveset {
	snakeMovesetLists := make([]Moveset, len(valueSnakeMap))
	totalCombinations := 1
	i := 0
	for _, snake := range valueSnakeMap {
		movesets := generatePossibleMoves(g, snake)
		snakeMovesetLists[i] = movesets
		totalCombinations *= len(movesets)
		i++
	}
	movesets := make([]Moveset, totalCombinations)
	combineMoves(movesets, snakeMoveLists)
}

func generateMovesets2(g game.Game)

/*func generateMovesets(g game.Game, valueSnakeMap game.SnakeByValue) <-chan Moveset {
	snakeMoveLists := make([][]game.Coordinate, len(valueSnakeMap))
	snakeOrder := make([]game.BoardValue, len(valueSnakeMap))
	totalCombinations := 1
	i := 0
	for snakeValue, snake := range valueSnakeMap {
		moveCoords := generatePossibleMoves(g, snake)
		snakeMoveLists[i] = moveCoords
		totalCombinations *= len(moveCoords)
		i++
	}
	movesets := make([]Moveset, totalCombinations)
	combineMoves(movesets, snakeMoveLists)
}*/

func generateOpponentGames(n node) []node {
	g := n.Game
	valueSnakeMap := g.ValueSnakeMap
	delete(valueSnakeMap, game.ME)
	moveSets := generateMovesets(g, valueSnakeMap)
	return nil
}
