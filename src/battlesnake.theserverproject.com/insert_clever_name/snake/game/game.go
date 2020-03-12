package game

var Games map[string]*Game

type Game struct {
	Id               string
	Board            Board
	PreviousMaxDepth int
}

type Board struct {
	Grid   Grid
	Snakes SnakeByValue
}

func CopyBoard(board Board) Board {
	return Board{
		Grid: copyGrid(board.Grid),
		Snakes: copySnakeByValues(board.Snakes),
	}
}

func InitGames() {
	Games = make(map[string]*Game)
}
