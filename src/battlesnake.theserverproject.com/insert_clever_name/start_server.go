package main

import (
	"battlesnake.theserverproject.com/insert_clever_name/snakeserver"
)

func main() {
	//Games := make(map[string]*game.Game)

	serverAddress := ":8001"
	server := snakeserver.NewSnakeServer(serverAddress, 5, 5)
	snakeserver.Start(server)
}
