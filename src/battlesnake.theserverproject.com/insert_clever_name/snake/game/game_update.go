package game

import (
	"errors"
)

type GameUpdate struct {
	Game struct {
		Id string `json:"id"`
	} `json:"game"`
	Turn  int `json:"turn"`
	Board struct {
		Height int           `json:"height"`
		Width  int           `json:"width"`
		Food   []Coordinate  `json:"food"`
		RawSnakes []snakeRaw `json:"snakes"`
	} `json:"board"`
	You snakeRaw `json:"you"`
}

func copyBoard(board Board) Board {
	height    := len(board)
	width     := len(board[0])
	boardCopy := make(Board, height)
	contents  := make([]BoardValue, height*width)
	for i := range board {
		start := i*width
		end   := start+width
		boardCopy[i] = contents[start:end:end]
		copy(boardCopy[i], board[i])
	}
	return boardCopy
}

func createBoard(state GameUpdate, snakesMap SnakeByValue) Board { // currently generates a new board every update for simplicity
	height   := state.Board.Height+2
	width    := state.Board.Width+2
	board    := make(Board, height)
	contents := make([]BoardValue, height*width)
	for i := range board {
		start := i*width
		end   := start+width
		board[i] = contents[start:end:end]
	}
	for _, snake := range snakesMap {
		for _, coordinate := range snake.Body {
			board[coordinate.Y][coordinate.X] = snake.Value
		}
	}
	for _, coordinate := range state.Board.Food {
		board[coordinate.Y][coordinate.X] = FOOD
	}
	return board
}

func CreateGame(state GameUpdate) {
	snakesMap := createSnakeMappings(state.Board.RawSnakes, state.You.Id)
	board     := createBoard(state, snakesMap)

	Games[state.Game.Id] = &Game{
		Id:            state.Game.Id,
		Board:         board,
		ValueSnakeMap: snakesMap,
	}
}

func UpdateGame(state GameUpdate) error {
	if game, ok := Games[state.Game.Id]; ok {
		game.ValueSnakeMap = createSnakeMappings(state.Board.RawSnakes, state.You.Id)
		game.Board         = createBoard(state, game.ValueSnakeMap)
		return nil
	} else {
		return errors.New("no game with given id for update")
	}
}

func DeleteGame(state GameUpdate) error {
	if _, ok := Games[state.Game.Id]; !ok {
		return errors.New("no game with given id for delete")
	}
	delete(Games, state.Game.Id) // garbage collector will do the rest
	return nil
}
