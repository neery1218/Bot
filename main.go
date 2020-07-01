package main

import (
	"automaton/automaton"
	"automaton/ofc"
	"log"
	"strconv"
	"time"
)

func Capture(imageName string, gCtxt ofc.GameContext) (*ofc.GameState, error) {
	err := gCtxt.CaptureGameState(imageName)
	if err != nil {
		return nil, err
	}

	gs, err := gCtxt.ParseGameStateFromImage(imageName)
	if err != nil {
		return nil, err
	}

	return gs, nil
}

func main() {
	log.Printf("Starting Bot\n")

	gCtxt := ofc.GameContext{
		Hwnd:             automaton.FindWindow("MEmu"),
		ImageServerHost:  "http://localhost:8000",
		SolverServerHost: "http://34.74.180.106:9001/eval",
		ScreenShotDir:    "C:\\Users\\neera\\Documents\\screenshots"}

	for {
		var gs *ofc.GameState
		var err error

		// take a screenshot every half second, break when we have to make a decision
		for {
			gs, err = Capture("tmp", gCtxt) // don't care about saving these
			if err != nil {                 // gs is nil
				log.Println("Line 42 failed.")
				log.Println(err)
			}

			if isValid, _ := gs.IsValid(); isValid && gs.DecisionRequired() {
				// sometimes we screenshot halfway through loading all the cards on the screen.
				// to fix that, sleep for half a second and recapture the game state
				time.Sleep(500 * time.Millisecond)

				timeString := strconv.Itoa(int(time.Now().Unix()))
				gs, err = Capture(timeString, gCtxt)
				if err != nil {
					log.Fatalf("Second screenshot %v failed %v", timeString, err)
				}

				if !gs.DecisionRequired() {
					panic("Somehow gamestate doesn't have a required decision")
				}
				break
			}
			log.Println("No Decision Required yet...")
			time.Sleep(200 * time.Millisecond)
		}

		log.Println("Decision needed!")
		actions, err := gCtxt.SolveGameState(gs)
		if err != nil {
			log.Println("Solver failed.")
			log.Fatal(err)
		}

		err = gCtxt.ExecuteActions(actions, gs)
		if err != nil {
			log.Fatalf("ExecuteActions failed: %+v\n", err)
		}

		// make sure actions were performed correctly and press the confirm button.
		err = gCtxt.ConfirmActions(actions)
		if err != nil {
			log.Println("Confirm Actions failed.")
			log.Fatal(err)
		}

		time.Sleep(200 * time.Millisecond)
	}
}
