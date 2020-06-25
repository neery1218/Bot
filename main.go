package main

import (
	"automaton/automaton"
	"automaton/ofc"
	"log"
	"strconv"
	"time"
)

func main() {
	log.Printf("Starting Bot\n")

	gCtxt := ofc.GameContext{
		Hwnd:             automaton.FindWindow("MEmu"),
		ImageServerHost:  "http://localhost:8000",
		SolverServerHost: "http://34.74.180.106:9001/eval",
		ScreenShotDir:    "C:\\Users\\neera\\Documents\\screenshots"}

	// screenshot id
	timeString := strconv.Itoa(int(time.Now().Unix()))

	err := gCtxt.CaptureGameState(timeString)
	if err != nil {
		panic(err)
	}

	gs, err := gCtxt.ParseGameStateFromImage(timeString)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", gs)
	if gs.DecisionRequired() {
		log.Println("Action required! Calling Solver")
		actions, err := gCtxt.SolveGameState(gs)
		if err != nil {
			panic(err)
		}
		log.Printf("Pending actions for timeString %v: %+v", timeString, actions)
		gCtxt.ExecuteActions(actions, gs)
	}
}
