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
	Hwnd             w32.HWND // memu handle
	ImageServerHost  string
	SolverServerHost string
	ScreenShotDir    string // store screenshots here
}

func (gCtxt *GameContext) CaptureGameState(imageName string) error {
	// get screenshot
	img := automaton.ScreenCapture(gCtxt.Hwnd)

	// save image to file
	f, err := os.Create(gCtxt.ScreenShotDir + "\\" + imageName + ".bmp")
	if err != nil {
		return err
	}

	// save image to file
	err = bmp.Encode(f, &img) // *Bitmap implements Image, not Bitmap!
	if err != nil {
		return err
	}

	return nil
}

func (gCtxt *GameContext) ParseGameStateFromImage(imageName string) (*GameState, error) {
	u, err := url.Parse(gCtxt.ImageServerHost)
	if err != nil {
		return nil, err
	}

	// image server needs a filename parameter, and config parameter
	q := u.Query()
	q.Set("filename", gCtxt.ScreenShotDir+"\\"+imageName+".bmp")
	q.Set("config", "MEMU_536x983")
	u.RawQuery = q.Encode()
	fmt.Printf("constructed url: %v\n", u)

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	gs, err := parseGameStateFromJson(string(body))
	if err != nil {
		return nil, err
	}
	return gs, nil
}

func (gCtxt *GameContext) SolveGameState(gs *GameState) error {
	return nil
}

func parseActions(actions string) {
	// "3d bot Ad mid 3c bot Ah mid Qs top"
}
