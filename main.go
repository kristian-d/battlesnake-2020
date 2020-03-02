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

	upChan := make(chan float64)
	downChan := make(chan float64)
	leftChan := make(chan float64)
	rightChan := make(chan float64)
	go ComputeMoveScore(upChan, Games[state.Game.Id], UP)
	go ComputeMoveScore(downChan, Games[state.Game.Id], DOWN)
	go ComputeMoveScore(leftChan, Games[state.Game.Id], LEFT)
	go ComputeMoveScore(rightChan, Games[state.Game.Id], RIGHT)

	var maxScore float64 = 0
	move := UP
	timeChan := time.NewTimer(time.Millisecond*250).C // process the move for 250ms, leaving approx. 250ms for the network
	for {
		select {
		case <-timeChan:
			res, err := json.Marshal(struct {
				Move string `json:"move"`
				Shout string `json:"shout"`
			}{
				Move: move,
				Shout: "shouting!",
			})
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
			return
		case upScore := <-upChan:
			if upScore > maxScore {
				move = UP
				maxScore = upScore
			}
		case downScore := <-downChan:
			if downScore > maxScore {
				move = DOWN
				maxScore = downScore
			}
		case leftScore := <-leftChan:
			if leftScore > maxScore {
				move = LEFT
				maxScore = leftScore
			}
		case rightScore := <-rightChan:
			if rightScore > maxScore {
				move = RIGHT
				maxScore = rightScore
			}
		}
	}
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
