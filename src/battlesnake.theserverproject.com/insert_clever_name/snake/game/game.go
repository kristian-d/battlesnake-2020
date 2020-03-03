package game

const (
	WALL  int = -1
	EMPTY int = 0
	FOOD  int = 1
)

var Games map[string]*Game

type Game struct {
	Id                  string
	Board               [][]int
	AliveSnakeCount     int
	SnakeValuesMap      map[string]*SnakeValues
	ValueSnakeValuesMap map[int]*SnakeValues
	Me                  *SnakeValues
}
