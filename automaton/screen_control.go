package automaton

import (
	"github.com/JamesHovious/w32"
)

func MoveMouse(hwnd w32.HWND, x, y int) {
	// translate relative coordinates to screen coordinates
	AbsX, AbsY := w32.ClientToScreen(hwnd, x, y)
	w32.SetCursorPos(AbsX, AbsY)
	return
}

func ClickDown() error {
	inputs := make([]w32.INPUT, 0)
	inputs = append(
		inputs,
		w32.INPUT{
			Type: w32.INPUT_MOUSE,
			Mi: w32.MOUSEINPUT{
				Dx:          0,
				Dy:          0,
				MouseData:   0,
				DwFlags:     w32.MOUSEEVENTF_LEFTDOWN,
				Time:        0,
				DwExtraInfo: 0},
			Ki: w32.KEYBDINPUT{},
			Hi: w32.HARDWAREINPUT{}})

	err := w32.SendInput(inputs)
	return err
}

func ClickUp() error {
	inputs := make([]w32.INPUT, 0)
	inputs = append(
		inputs,
		w32.INPUT{
			Type: w32.INPUT_MOUSE,
			Mi: w32.MOUSEINPUT{
				Dx:          0,
				Dy:          0,
				MouseData:   0,
				DwFlags:     w32.MOUSEEVENTF_LEFTUP,
				Time:        0,
				DwExtraInfo: 0},
			Ki: w32.KEYBDINPUT{},
			Hi: w32.HARDWAREINPUT{}})

	err := w32.SendInput(inputs)
	return err
}
