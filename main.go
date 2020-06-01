package main

import (
	"automaton/automaton"
	"automaton/ofc"
	"fmt"
)

func main() {

	gCtxt := ofc.GameContext{
		Hwnd: automaton.FindWindow("MEmu"),
		Host: "http://localhost:8000"}

	gs := ofc.CaptureGameState(gCtxt)
	fmt.Println(gs)
}
