package game

import "fmt"

type BoardValue uint8
const (
	EMPTY BoardValue = iota
	FOOD
	ME
)

type Board [][]BoardValue

func createBoard(state GameUpdate, snakesMap SnakeByValue) Board { // currently generates a new board every update for simplicity
	height   := state.Board.Height
	width    := state.Board.Width
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

func PrintBoard(board Board) {
	for _, row := range board {
		fmt.Println(row)
	}
}
