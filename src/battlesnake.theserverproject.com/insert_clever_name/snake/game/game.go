package game

var Games map[string]*Game

type Game struct {
	Id            string
	Board         Board
	ValueSnakeMap SnakeByValue
}

func CopyGame(game Game) Game {
	return Game{
		Id:            game.Id,
		Board:         copyBoard(game.Board),
		ValueSnakeMap: copySnakeByValues(game.ValueSnakeMap),
	}
}

func InitGames() {
	Games = make(map[string]*Game)
}
