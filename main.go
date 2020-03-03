package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "sssSssSsssssSs I am alive! sssSssSsssssSs")
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "sssSssSsssssSs I am alive! sssSssSsssssSs")
}

func clean(req *http.Request) (*GameUpdate, error) {
	info := &GameUpdate{}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(b, info)
	return info, err
}

func start(w http.ResponseWriter, req *http.Request) {
	state, err := clean(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	CreateGame(state)
	fmt.Printf("Game created, id=%s\n", state.Game.Id)

	res, err := json.Marshal(struct {
		Color string `json:"color"`
		HeadType string `json:"headType"`
		TailType string `json:"tailType"`
	}{
		Color: "#ff00ff",
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

func move(w http.ResponseWriter, req *http.Request) {
	state, err := clean(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = UpdateGame(state)
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
		Move:  ComputeMove(Games[state.Game.Id], 250), // process the move for x ms, leaving (500 - x) ms for the network
		Shout: "shouting!",
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func end(w http.ResponseWriter, req *http.Request) {
	state, err := clean(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = DeleteGame(state)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Printf("Cannot delete game, id=%s\n", state.Game.Id)
		return
	}
	fmt.Printf("Game deleted, id=%s\n", state.Game.Id)
}

func main() {
	Games = make(map[string]*Game)

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/start", start)
	mux.HandleFunc("/move", move)
	mux.HandleFunc("/end", end)

	address := ":8081"
	srv := &http.Server{
		Addr: address,
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	fmt.Printf("Server listening on port %s\n", address[1:])
	log.Fatal(srv.ListenAndServe())
}
