package snakeserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"battlesnake.theserverproject.com/insert_clever_name/snake/engine"
	"battlesnake.theserverproject.com/insert_clever_name/snake/game"
)

// SnakeServer An http server used to serve
// moves played by the snake.
type SnakeServer struct {
	serverAddress string
	readTimeout   int
	writeTimeout  int
}

// NewSnakeServer Instanciates a new snake server.
func NewSnakeServer(serverAddress string, readTimeout int, writeTimeout int) SnakeServer {
	return SnakeServer{serverAddress, readTimeout, writeTimeout}
}

// Start starts the snake server using mux
func Start(server SnakeServer) {
	game.InitGames()
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/start", start)
	mux.HandleFunc("/move", move)
	mux.HandleFunc("/end", end)

	srv := &http.Server{
		Addr:         server.serverAddress,
		Handler:      mux,
		ReadTimeout:  time.Duration(server.readTimeout) * time.Second,
		WriteTimeout: time.Duration(server.writeTimeout) * time.Second,
	}

	fmt.Printf("Server listening on address %s\n", server.serverAddress)
	log.Fatal(srv.ListenAndServe())
}

func clean(req *http.Request) (game.GameUpdate, error) {
	info := &game.GameUpdate{}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return *info, err
	}

	err = json.Unmarshal(b, info)
	return *info, err
}

func end(w http.ResponseWriter, req *http.Request) {
	state, err := clean(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = game.DeleteGame(state)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Printf("Cannot delete game, id=%s\n", state.Game.Id)
		return
	}
	fmt.Printf("Game deleted, id=%s\n", state.Game.Id)
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "sssSssSsssssSs I am alive! sssSssSsssssSs")
}

func move(w http.ResponseWriter, req *http.Request) {
	state, err := clean(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = game.UpdateGame(state)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Printf("Cannot update game, id=%s\n", state.Game.Id)
		return
	}
	fmt.Printf("Game updated, id=%s\n", state.Game.Id)

	res, err := json.Marshal(struct {
		Move  string `json:"move"`
		Shout string `json:"shout"`
	}{
		Move:  engine.ComputeMove(*game.Games[state.Game.Id], 250), // process the move for x ms, leaving (500 - x) ms for the network
		Shout: "shouting!",
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "sssSssSsssssSs I am alive! sssSssSsssssSs")
}

func start(w http.ResponseWriter, req *http.Request) {
	state, err := clean(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	game.CreateGame(state)
	fmt.Printf("Game created, id=%s\n", state.Game.Id)

	res, err := json.Marshal(struct {
		Color    string `json:"color"`
		HeadType string `json:"headType"`
		TailType string `json:"tailType"`
	}{
		Color:    "#ff00ff",
		HeadType: "bendr",
		TailType: "pixel",
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
