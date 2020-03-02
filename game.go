package main

import "errors"

const (
	WALL int = -1
	EMPTY int = 0
	FOOD int = 1
)

var Games map[string]*Game

func createSnakeMappings(snakes []Snake) (map[string]*SnakeValues, map[int]*SnakeValues) {
	snakeValuesMap := make(map[string]*SnakeValues)
	valueSnakeValuesMap := make(map[int]*SnakeValues)
	for i, snake := range snakes { // TODO: does the Board's snake list already include myself?
		snakeValues := SnakeValues{
			Id: snake.Id,
			Size: len(snake.Body),
			Health: snake.Health,
			HeadCoordinate: snake.Body[0],
			TailCoordinate: snake.Body[len(snake.Body) - 1],
			HeadValue: i*3 + 1 + FOOD, // ensures that values don't interfere with FOOD or one another
			BodyValue: i*3 + 2 + FOOD,
			TailValue: i*3 + 3 + FOOD,
		}
		snakeValuesMap[snake.Id] = &snakeValues
		valueSnakeValuesMap[snakeValues.HeadValue] = &snakeValues
		valueSnakeValuesMap[snakeValues.BodyValue] = &snakeValues
		valueSnakeValuesMap[snakeValues.TailValue] = &snakeValues
	}
	return snakeValuesMap, valueSnakeValuesMap
}

func updateBoard(state *GameUpdate) { // currently generates a new board every update for simplicity
	board := make([][]int, state.Board.Height + 2)
	for i := range board {
		board[i] = make([]int, state.Board.Width + 2)
	}
	snakeValuesMap := Games[state.Game.Id].SnakeValuesMap
	allSnakes := append(state.Board.Snakes, state.You)
	for _, snake := range allSnakes { // TODO: does the Board's snake list already include myself?
		for j, coordinate := range snake.Body {
			if j == 0 {
				board[coordinate.Y][coordinate.X] = snakeValuesMap[snake.Id].HeadValue
			} else if j == len(snake.Body) - 1 {
				board[coordinate.Y][coordinate.X] = snakeValuesMap[snake.Id].TailValue
			} else {
				board[coordinate.Y][coordinate.X] = snakeValuesMap[snake.Id].BodyValue
			}
		}
	}
	for _, coordinate := range state.Board.Food {
		board[coordinate.Y][coordinate.X] = FOOD
	}
	Games[state.Game.Id].Board = board
}

func updateSnakeMappings(state *GameUpdate) {
	allSnakes := append(state.Board.Snakes, state.You)
	for _, snake := range allSnakes { // TODO: does the Board's snake list already include myself?
		snakeValuesMap := Games[state.Game.Id].SnakeValuesMap[snake.Id]
		snakeValuesMap.Health = snake.Health
		snakeValuesMap.Size = len(snake.Body)
		snakeValuesMap.HeadCoordinate = snake.Body[0]
		snakeValuesMap.TailCoordinate = snake.Body[len(snake.Body) - 1]
	}
}

func CreateGame(state *GameUpdate) {
	board := make([][]int, state.Board.Height + 2)
	for i := range board {
		board[i] = make([]int, state.Board.Width + 2)
	}

	allSnakes := append(state.Board.Snakes, state.You)
	snakeValuesMap, valueSnakeValuesMap := createSnakeMappings(allSnakes)

	Games[state.Game.Id] = &Game{
		Id: state.Game.Id,
		Board: board,
		AliveSnakeCount: len(state.Board.Snakes) + 1, // TODO: does the Board's snake list include myself?
		SnakeValuesMap: snakeValuesMap,
		ValueSnakeValuesMap: valueSnakeValuesMap,
		Me: snakeValuesMap[state.You.Id],
	}

	updateBoard(state)
}

func UpdateGame(state *GameUpdate) error {
	if _, ok := Games[state.Game.Id]; !ok {
		return errors.New("no game with given id for update")
	}
	updateBoard(state)
	updateSnakeMappings(state)
	return nil
}

func DeleteGame(state *GameUpdate) error {
	if _, ok := Games[state.Game.Id]; !ok {
		return errors.New("no game with given id for delete")
	}
	delete(Games, state.Game.Id) // garbage collector will do the rest
	return nil
}
