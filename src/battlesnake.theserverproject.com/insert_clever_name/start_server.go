package main

import (
	"battlesnake.theserverproject.com/insert_clever_name/snakeserver"
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("Detected %d CPUS\n", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("Utilizing %d CPUS\n", runtime.NumCPU())
	serverAddress := ":8001"
	server := snakeserver.NewSnakeServer(serverAddress, 5, 5)
	snakeserver.Start(server)
}
