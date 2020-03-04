package game

const (
	EMPTY = iota
	WALL
	FOOD
	ME
)

var Games map[string]*Game

type Game struct {
	Id                  string
	Board               [][]int
	ValueSnakeMap       map[int]Snake
}

func InitGames() {
	Games = make(map[string]*Game)
}
