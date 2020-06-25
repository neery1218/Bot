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

	for {
		var gs *ofc.GameState
		var timeString string
		var err error

		// take a screenshot every half second, break when we have to make a decision
		for {
			timeString = strconv.Itoa(int(time.Now().Unix()))

			err = gCtxt.CaptureGameState(timeString)
			if err != nil {
				panic(err)
			}

			gs, err = gCtxt.ParseGameStateFromImage(timeString)
			if err != nil {
				panic(err)
			}
			if gs.DecisionRequired() {
				break
			}
			log.Println("No Decision Required yet...")
			time.Sleep(500 * time.Millisecond)
		}

		if !gs.DecisionRequired() {
			panic("Somehow gamestate doesn't have a required decision")
		}
		log.Println("Decision needed!")
		actions, err := gCtxt.SolveGameState(gs)
		if err != nil {
			panic(err)
		}
		gCtxt.ExecuteActions(actions, gs)
		time.Sleep(1000 * time.Millisecond)
	}
}
