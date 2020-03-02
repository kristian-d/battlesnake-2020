package main

const (
	UP string = "up"
	DOWN string = "down"
	LEFT string = "left"
	RIGHT string = "right"
)

func copyBoard(board *[][]int) *[][]int {
	n := len(*board)
	m := len((*board)[0])
	duplicate := make([][]int, n)
	data := make([]int, n*m)
	for i := range *board {
		start := i*m
		end := start + m
		duplicate[i] = data[start:end:end]
		copy(duplicate[i], (*board)[i])
	}
	return &duplicate
}

func moveSnake(board *[][]int, me SnakeValues, move string) {
	switch move {
	case UP:
		newHeadCoordinate :=
	}
}

func BFS() {

}

func ComputeMoveScore(scoreChan chan<-float64, game *Game, move string) {
	board := copyBoard(&game.Board)
	moveSnake(board, move)

	// calculate the amount of space each other snake has by doing a breadth first search for each other snake
}
