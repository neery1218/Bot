package main

import (
	"automaton/automaton"
	"os"

	"golang.org/x/image/bmp"
)

func main() {

	// Notepad has to be open or this won't work.
	appName := "Notepad"
	img := automaton.ScreenCapture(appName)
	f, err := os.Create("notepad.bmp")
	if err != nil {
		panic(err)
	}
	bmp.Encode(f, &img) // *Bitmap implements Image, not Bitmap!

	automaton.MoveMouse(appName, 10, 11)
	automaton.ClickDown(appName)
	automaton.MoveMouse(appName, 20, 11)
	// automaton.ClickUp(appName)
}
