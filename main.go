package main

import (
	"automaton/automaton"
	"automaton/ofc"
	"fmt"
	"strconv"
	"time"
)

func main() {

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
		fmt.Println("Action required! Calling Solver")
		actions, err := gCtxt.SolveGameState(gs)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Pending actions: %+v", actions)
		gCtxt.ExecuteActions(actions, gs)
	}
}
