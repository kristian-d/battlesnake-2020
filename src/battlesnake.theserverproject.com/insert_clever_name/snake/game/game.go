package game

var Games map[string]*Game

type Game struct {
	Id            string
	Grid          Grid
	ValueSnakeMap SnakeByValue
}

func CopyGame(game Game) Game {
	return Game{
		Id:            game.Id,
		Grid:          copyGrid(game.Grid),
		ValueSnakeMap: copySnakeByValues(game.ValueSnakeMap),
	}
}

func InitGames() {
	Games = make(map[string]*Game)
}
