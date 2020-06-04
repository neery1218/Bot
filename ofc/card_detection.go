package ofc

import (
	"automaton/automaton"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/JamesHovious/w32"
	"golang.org/x/image/bmp"
)

type GameContext struct {
	Hwnd w32.HWND
	Host string
}

func CaptureGameState(ctxt GameContext) *GameState {
	// get screenshot
	img := automaton.ScreenCapture(ctxt.Hwnd)

	// save image to file
	filename := "C:\\Users\\neera\\Documents\\screenshots\\tmp.bmp"
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	bmp.Encode(f, &img) // *Bitmap implements Image, not Bitmap!

	// call card detection server
	u, err := url.Parse(ctxt.Host)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("filename", filename)
	q.Set("config", "MEMU_536x983")
	u.RawQuery = q.Encode()
	fmt.Printf("constructed url: %v\n", u)

	res, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	gs, err := parseGameStateFromJson(string(body))
	if err != nil {
		panic(err)
	}
	return gs
}
