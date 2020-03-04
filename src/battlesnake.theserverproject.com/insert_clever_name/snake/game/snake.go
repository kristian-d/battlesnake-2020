package game

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SnakeRaw struct {
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
	Value          int
	Alive          bool
}

func createSnakeMappings(rawSnakes []SnakeRaw, myId string) map[int]Snake {
	snakesMapping := make(map[int]Snake)
	for i, rawSnake := range rawSnakes {
		value := i + 1 + ME // ensures that values are unique
		if rawSnake.Id == myId {
			value = ME
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
