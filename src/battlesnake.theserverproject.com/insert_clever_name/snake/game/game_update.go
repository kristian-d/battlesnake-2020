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

func CreateGame(state GameUpdate) {
	snakesMap := createSnakeMappings(state.Board.RawSnakes, state.You.Id)
	grid      := createGrid(state, snakesMap)

	Games[state.Game.Id] = &Game{
		Id:            state.Game.Id,
		Grid:          grid,
		ValueSnakeMap: snakesMap,
	}
}

func UpdateGame(state GameUpdate) error {
	if game, ok := Games[state.Game.Id]; ok {
		game.ValueSnakeMap = createSnakeMappings(state.Board.RawSnakes, state.You.Id)
		game.Grid          = createGrid(state, game.ValueSnakeMap)
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
