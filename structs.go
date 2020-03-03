package main

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Snake struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Health int `json:"health"`
	Body []Coordinate `json:"body"`
	Shout string `json:"shout"`
}

type GameUpdate struct {
	Game struct {
		Id string `json:"id"`
	} `json:"game"`
	Turn int `json:"turn"`
	Board struct {
		Height int `json:"height"`
		Width int `json:"width"`
		Food []Coordinate `json:"food"`
		Snakes []Snake `json:"snakes"`
	} `json:"board"`
	You Snake `json:"you"`
}

type SnakeValues struct {
	Id string
	Size int
	Health int
	HeadCoordinate Coordinate
	TailCoordinate Coordinate
	HeadValue int
	BodyValue int
	TailValue int
}

type Game struct {
	Id string
	Board [][]int
	AliveSnakeCount int
	SnakeValuesMap map[string]*SnakeValues
	ValueSnakeValuesMap map[int]*SnakeValues
	Me *SnakeValues
}

type Node struct {
	Board [][]int
	Children []Node
	Expanded bool
}
