package main

import (
	"automaton/automaton"
	"os"

	"golang.org/x/image/bmp"
)

func main() {

	// Notepad has to be open or this won't work.
	img := automaton.ScreenCapture("Notepad")
	f, err := os.Create("notepad.bmp")
	if err != nil {
		panic(err)
	}
	bmp.Encode(f, &img) // *Bitmap implements Image, not Bitmap!
}
