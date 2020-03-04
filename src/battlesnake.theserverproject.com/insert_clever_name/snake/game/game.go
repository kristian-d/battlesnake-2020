package game

const (
	WALL    int = -1
	EMPTY   int = 0
	FOOD    int = 1
	ME      int = 2
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
