package game

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Id     string       `json:"id"`
	Name   string       `json:"name"`
	Health int          `json:"health"`
	Body   []Coordinate `json:"body"`
	Shout  string       `json:"shout"`
}

type SnakeValues struct {
	Id             string
	Size           int
	Health         int
	HeadCoordinate Coordinate
	TailCoordinate Coordinate
	HeadValue      int
	BodyValue      int
	TailValue      int
}

func createSnakeMappings(snakes []Snake) (map[string]*SnakeValues, map[int]*SnakeValues) {
	snakeValuesMap := make(map[string]*SnakeValues)
	valueSnakeValuesMap := make(map[int]*SnakeValues)
	for i, snake := range snakes { // TODO: does the Board's snake list already include myself?
		snakeValues := SnakeValues{
			Id:             snake.Id,
			Size:           len(snake.Body),
			Health:         snake.Health,
			HeadCoordinate: snake.Body[0],
			TailCoordinate: snake.Body[len(snake.Body)-1],
			HeadValue:      i*3 + 1 + FOOD, // ensures that values don't interfere with FOOD or one another
			BodyValue:      i*3 + 2 + FOOD,
			TailValue:      i*3 + 3 + FOOD,
		}
		snakeValuesMap[snake.Id] = &snakeValues
		valueSnakeValuesMap[snakeValues.HeadValue] = &snakeValues
		valueSnakeValuesMap[snakeValues.BodyValue] = &snakeValues
		valueSnakeValuesMap[snakeValues.TailValue] = &snakeValues
	}
	return snakeValuesMap, valueSnakeValuesMap
}
