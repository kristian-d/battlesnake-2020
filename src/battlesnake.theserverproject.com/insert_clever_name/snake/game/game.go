package game

type BoardValue int
const (
	EMPTY BoardValue = iota
	WALL
	FOOD
	ME
)

var Games map[string]*Game

type Board [][]BoardValue
type SnakeByValue map[BoardValue]Snake
type Game struct {
	Id            string
	Board         Board
	ValueSnakeMap SnakeByValue
}

func InitGames() {
	Games = make(map[string]*Game)
}
