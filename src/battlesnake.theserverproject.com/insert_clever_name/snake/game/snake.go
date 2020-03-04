package game

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type snakeRaw struct {
	Id     string       `json:"id"`
	Name   string       `json:"name"`
	Health int          `json:"health"`
	Body   []Coordinate `json:"body"`
	Shout  string       `json:"shout"`
}

type Snake struct {
	Id             string
	Health         int
	Body           []Coordinate
	Value          BoardValue
	Alive          bool
}

func createSnakeMappings(rawSnakes []snakeRaw, myId string) map[BoardValue]Snake {
	snakesMapping := make(map[BoardValue]Snake)
	for i, rawSnake := range rawSnakes {
		var value BoardValue
		if rawSnake.Id == myId {
			value = ME
		} else {
			value = BoardValue(i + 1) + ME // ensures that values are unique
		}
		snakesMapping[value] = Snake{
			Id:             rawSnake.Id,
			Health:         rawSnake.Health,
			Body:           rawSnake.Body,
			Value:          value,
			Alive:          true,
		}
	}
	return snakesMapping
}
